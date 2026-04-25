package monitoring

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
	wbmapper "salehelp/internal/adapter/wb/mapper"
	"salehelp/internal/storage/pggorm/model"
	_ "salehelp/internal/storage/pggorm/repository/shop"
)

// базовый URL для формирования ссылки на товар WB
const wbProductBaseURL = "https://www.wildberries.ru"

const (
	// heavySyncInterval — минимальный интервал между тяжёлыми синками (карточки, склады, снапшоты).
	// Контент API WB — лимит 1 запрос/час на базовом токене.
	heavySyncInterval = time.Hour

	// fullCardsSyncInterval — интервал полного (не инкрементального) синка карточек.
	// Нужен для выявления удалённых товаров, которые инкрементальный синк не увидит.
	fullCardsSyncInterval = 7 * 24 * time.Hour
)

func (s *Service) SyncAllShops(ctx context.Context) error {
	shops, err := s.shopsRepo.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("sync all shops: list shops: %w", err)
	}

	for _, shop := range shops {
		if err := s.SyncShop(ctx, shop.ID); err != nil {
			s.log.Error("sync shop failed", "shop_id", shop.ID, "error", err)
		}
	}

	return nil
}

func (s *Service) SyncShop(ctx context.Context, shopID int) (err error) {
	shop, err := s.shopsRepo.GetByID(ctx, shopID)
	if err != nil {
		return fmt.Errorf("sync shop: get shop by id %d: %w", shopID, err)
	}
	if shop == nil {
		return fmt.Errorf("sync shop: shop %d not found", shopID)
	}
	if shop.APIToken == "" {
		return fmt.Errorf("sync shop: shop %d has empty WB API key", shopID)
	}

	defer func() {
		now := time.Now()
		upd := model.SyncStatusUpdate{LastSyncAt: &now}

		if err == nil {
			upd.LastSyncOkAt = &now
		} else {
			upd.LastSyncFailedAt = &now
			if sc, ok := wbStatusCode(err); ok {
				upd.LastSyncStatusCode = &sc
			}
			m := err.Error()
			upd.LastSyncError = &m
			s.log.Error("sync_shop: critical failure — data collection stopped",
				"shop_id", shopID, "error", err.Error(), "critical", true,
			)
		}

		_ = s.shopsRepo.UpdateSyncStatus(ctx, shopID, upd)
	}()

	if err = s.runFastSync(ctx, shop); err != nil {
		return fmt.Errorf("sync shop: fast sync: %w", err)
	}

	if s.shouldRunHeavySync(shop) {
		if err = s.runHeavySync(ctx, shop); err != nil {
			return fmt.Errorf("sync shop: heavy sync: %w", err)
		}
	}

	if shop.FirstSyncNotifiedAt == nil {
		n, err2 := s.productsRepo.CountByShop(ctx, shopID)
		if err2 == nil {
			user, _ := s.userRepo.GetByID(ctx, shop.UserID)
			chatID := user.TelegramChatID
			if chatID != nil && *chatID != 0 {
				msg := fmt.Sprintf(
					"✅ Первый синк завершен!\n\n"+
						"Я нашел в вашем магазине %d товаров.\n"+
						"Теперь я слежу за ними 24/7. Как только что-то изменится — я сразу маякну!",
					n,
				)
				if err3 := s.tg.SendMessage(ctx, *chatID, msg); err3 == nil {
					_ = s.shopsRepo.UpdateFirstSyncNotifiedAt(ctx, shopID, time.Now())
				}
			}
		}
	}

	return nil
}

func (s *Service) shouldRunHeavySync(shop *domain.Shop) bool {
	if shop.LastHeavySyncAt == nil {
		return true
	}
	return time.Since(*shop.LastHeavySyncAt) >= heavySyncInterval
}

// runFastSync выполняется при каждом тике планировщика (каждые 30 мин).
// Содержит только вызовы с высокими лимитами WB: заказы, отзывы, продажи.
func (s *Service) runFastSync(ctx context.Context, shop *domain.Shop) error {
	if err := s.syncShopFBSOrders(ctx, shop); err != nil {
		return fmt.Errorf("fbs orders: %w", err)
	}
	if err := s.syncShopDBWOrders(ctx, shop); err != nil {
		return fmt.Errorf("dbw orders: %w", err)
	}
	if err := s.syncShopDBSOrders(ctx, shop); err != nil {
		return fmt.Errorf("dbs orders: %w", err)
	}
	if err := s.syncShopFeedbacks(ctx, shop); err != nil {
		return fmt.Errorf("feedbacks: %w", err)
	}
	if err := s.syncShopQuestions(ctx, shop); err != nil {
		return fmt.Errorf("questions: %w", err)
	}
	if err := s.syncShopSales(ctx, shop); err != nil {
		return fmt.Errorf("sales: %w", err)
	}
	if err := s.syncShopOrderRefusals(ctx, shop); err != nil {
		return fmt.Errorf("order refusals: %w", err)
	}
	return nil
}

// runHeavySync выполняется не чаще 1 раза в час (ограничен лимитами Content API WB).
// Синхронизирует склады, карточки товаров (инкрементально) и снапшоты цен/остатков.
func (s *Service) runHeavySync(ctx context.Context, shop *domain.Shop) error {
	if err := s.syncShopWarehouses(ctx, shop); err != nil {
		return fmt.Errorf("warehouses: %w", err)
	}

	newCursor, isFullSync, err := s.syncShopProductsIncremental(ctx, shop)
	if err != nil {
		return fmt.Errorf("products: %w", err)
	}

	// снапшоты делаем по всем товарам магазина, а не только по изменённым
	allPtrs, err := s.productsRepo.ListByShop(ctx, shop.ID, 0, 0)
	if err != nil {
		return fmt.Errorf("list products for snapshots: %w", err)
	}

	if len(allPtrs) > 0 {
		products := make([]domain.Product, len(allPtrs))
		for i, p := range allPtrs {
			products[i] = *p
		}
		if err := s.syncShopSnapshots(ctx, shop, products); err != nil {
			return fmt.Errorf("snapshots: %w", err)
		}
	}

	now := time.Now()
	upd := model.HeavySyncUpdate{LastHeavySyncAt: &now}

	if newCursor != nil {
		upd.CardsCursorUpdatedAt = &newCursor.UpdatedAt
		nmID := newCursor.NmID
		upd.CardsCursorNmID = &nmID
	}
	if isFullSync {
		upd.LastFullCardsSyncAt = &now
	}

	if err := s.shopsRepo.UpdateHeavySyncState(ctx, shop.ID, upd); err != nil {
		return fmt.Errorf("update heavy sync state: %w", err)
	}

	return nil
}

// syncShopFBSOrders получает новые сборочные задания FBS из WB и сохраняет в fbs_orders.
// Дубликаты (shop_id, wb_order_id) игнорируются.
func (s *Service) syncShopFBSOrders(ctx context.Context, shop *domain.Shop) error {
	resp, err := s.wbClient.GetNewOrders(ctx, shop.APIToken)
	if err != nil {
		return fmt.Errorf("wb get new orders: %w", err)
	}
	if len(resp.Orders) == 0 {
		return nil
	}
	for i := range resp.Orders {
		o := wbmapper.FBSOrderNewToDomain(resp.Orders[i], shop.ID)
		if err := s.fbsOrdersRepo.Create(ctx, &o); err != nil {
			return fmt.Errorf("create fbs_order wb_order_id=%d: %w", o.WBOrderID, err)
		}
	}
	return nil
}

// syncShopDBWOrders получает новые сборочные задания DBW из WB и сохраняет в dbw_orders.
// Дубликаты (shop_id, wb_order_id) игнорируются. Используется RateMarketplace (rate limiter).
func (s *Service) syncShopDBWOrders(ctx context.Context, shop *domain.Shop) error {
	resp, err := s.wbClient.GetNewDBWOrders(ctx, shop.APIToken)
	if err != nil {
		return fmt.Errorf("wb get new dbw orders: %w", err)
	}
	if len(resp.Orders) == 0 {
		return nil
	}
	for i := range resp.Orders {
		o := wbmapper.DBWOrderNewToDomain(resp.Orders[i], shop.ID)
		if err := s.dbwOrdersRepo.Create(ctx, &o); err != nil {
			return fmt.Errorf("create dbw_order wb_order_id=%d: %w", o.WBOrderID, err)
		}
	}
	return nil
}

// syncShopDBSOrders получает новые сборочные задания DBS из WB и сохраняет в dbs_orders.
// Дубликаты (shop_id, wb_order_id) игнорируются. Используется RateMarketplace (rate limiter).
func (s *Service) syncShopDBSOrders(ctx context.Context, shop *domain.Shop) error {
	resp, err := s.wbClient.GetNewDBSOrders(ctx, shop.APIToken)
	if err != nil {
		return fmt.Errorf("wb get new dbs orders: %w", err)
	}
	if len(resp.Orders) == 0 {
		return nil
	}
	for i := range resp.Orders {
		o := wbmapper.DBSOrderNewToDomain(resp.Orders[i], shop.ID)
		if err := s.dbsOrdersRepo.Create(ctx, &o); err != nil {
			return fmt.Errorf("create dbs_order wb_order_id=%d: %w", o.WBOrderID, err)
		}
	}
	return nil
}

// syncShopFeedbacks получает необработанные отзывы WB (GET /api/v1/feedbacks) и сохраняет в wb_feedbacks.
// Лимит: 3 запроса/сек. Пагинация take=1000.
func (s *Service) syncShopFeedbacks(ctx context.Context, shop *domain.Shop) error {
	const pageSize = 1000
	skip := 0
	for {
		resp, err := s.wbClient.GetFeedbacks(ctx, shop.APIToken, false, pageSize, skip, "dateDesc", nil, nil)
		if err != nil {
			return fmt.Errorf("wb get feedbacks: %w", err)
		}
		if len(resp.Data.Feedbacks) == 0 {
			break
		}
		for i := range resp.Data.Feedbacks {
			f := wbmapper.WBFeedbackToDomain(resp.Data.Feedbacks[i], shop.ID)
			if err := s.wbFeedbacksRepo.Create(ctx, &f); err != nil {
				return fmt.Errorf("create wb_feedback id=%s: %w", f.WBFeedbackID, err)
			}
		}
		if len(resp.Data.Feedbacks) < pageSize {
			break
		}
		skip += pageSize
	}
	return nil
}

// syncShopQuestions получает неотвеченные вопросы WB (GET /api/v1/questions) и сохраняет в wb_questions.
// Лимит: 3 запроса/сек. Пагинация take=1000.
func (s *Service) syncShopQuestions(ctx context.Context, shop *domain.Shop) error {
	const pageSize = 1000
	skip := 0
	for {
		resp, err := s.wbClient.GetQuestions(ctx, shop.APIToken, false, pageSize, skip, "dateDesc", nil, nil)
		if err != nil {
			return fmt.Errorf("wb get questions: %w", err)
		}
		if len(resp.Data.Questions) == 0 {
			break
		}
		for i := range resp.Data.Questions {
			q := wbmapper.WBQuestionToDomain(resp.Data.Questions[i], shop.ID)
			if err := s.wbQuestionsRepo.Create(ctx, &q); err != nil {
				return fmt.Errorf("create wb_question id=%s: %w", q.WBQuestionID, err)
			}
		}
		if len(resp.Data.Questions) < pageSize {
			break
		}
		skip += pageSize
	}
	return nil
}

// syncShopSales получает отчёт продаж/возвратов WB (GET /api/v1/supplier/sales) и сохраняет в sale_events.
// Дубликаты (shop_id, srid) игнорируются. Пагинация по lastChangeDate (лимит WB ~80k строк за запрос).
// При первом синке (нет записей в sale_events) загружаем последние 90 дней, помечаем как notified,
// отправляем сводку (не отдельные выкупы). Со второго синка — уведомления в реальном времени.
func (s *Service) syncShopSales(ctx context.Context, shop *domain.Shop) error {
	// WB отдаёт даты в отчёте продаж без таймзоны (bare datetime) в московском времени.
	// Парсим через глобальный wbReportLoc (Europe/Moscow из конфига).
	// Для отображения в уведомлениях используется shopLoc из таблицы shops (см. alert/engine.go).
	dateFrom, err := s.saleEventsRepo.GetMaxLastChangeDateByShop(ctx, shop.ID)
	if err != nil {
		return fmt.Errorf("get max last_change_date: %w", err)
	}
	isFirstSync := dateFrom.IsZero()
	if isFirstSync {
		dateFrom = time.Now().AddDate(0, 0, -90)
	} else {
		// Запрашиваем строго после последней записи, чтобы не зависеть от трактовки >= на стороне WB
		dateFrom = dateFrom.Add(time.Millisecond)
	}

	var buyouts, returns int
	var revenue float64

	// При первом синке сразу вставляем с notified_at, чтобы alert worker не успел разослать до MarkAllNotifiedByShop
	var notifiedAt *time.Time
	if isFirstSync {
		now := time.Now()
		notifiedAt = &now
	}

	var batch []*domain.SaleEvent
	for {
		resp, err := s.wbClient.GetSales(ctx, shop.APIToken, dateFrom, 0, nil)
		if err != nil {
			return fmt.Errorf("wb get sales: %w", err)
		}
		if len(resp) == 0 {
			break
		}
		var maxLastChange time.Time
		for i := range resp {
			item := resp[i]
			e := wbmapper.SaleItemToDomain(item, shop.ID, nil)
			e.NotifiedAt = notifiedAt
			ep := new(domain.SaleEvent)
			*ep = e
			batch = append(batch, ep)
			if item.ForPay >= 0 {
				buyouts++
				revenue += item.ForPay
			} else {
				returns++
				revenue += item.ForPay
			}
			if lc, err := item.ParseLastChangeDate(nil); err == nil && lc.After(maxLastChange) {
				maxLastChange = lc
			}
		}
		if !maxLastChange.IsZero() {
			dateFrom = maxLastChange.Add(time.Millisecond)
		} else if len(resp) > 0 {
			if lastChange, _ := resp[len(resp)-1].ParseLastChangeDate(nil); !lastChange.IsZero() {
				dateFrom = lastChange.Add(time.Millisecond)
			}
		}
	}
	if len(batch) > 0 {
		if err := s.saleEventsRepo.CreateBatch(ctx, batch); err != nil {
			return fmt.Errorf("create sale_events batch: %w", err)
		}
	}

	if isFirstSync {
		now := time.Now()
		if err := s.saleEventsRepo.MarkAllNotifiedByShop(ctx, shop.ID, now); err != nil {
			return fmt.Errorf("mark sale_events notified: %w", err)
		}
		msg := formatFirstSyncSalesSummaryRU(buyouts, returns, revenue)
		if msg != "" {
			chatIDs, err := s.endpointsRepo.ListEnabledTelegramChatIDsByShopID(ctx, shop.ID)
			if err != nil {
				s.log.Warn("first sync sales summary: list endpoints failed", "shop_id", shop.ID, "err", err)
			}
			if len(chatIDs) == 0 {
				user, _ := s.userRepo.GetByID(ctx, shop.UserID)
				if user != nil && user.TelegramChatID != nil && *user.TelegramChatID != 0 {
					chatIDs = []int64{*user.TelegramChatID}
				}
			}
			for _, chatID := range chatIDs {
				if err := s.tg.SendMessage(ctx, chatID, msg); err != nil {
					s.log.Warn("first sync sales summary: send failed", "shop_id", shop.ID, "chat_id", chatID, "err", err)
				}
			}
		}
	}

	return nil
}

// syncShopOrderRefusals получает заказы WB (GET /api/v1/supplier/orders), сохраняет отказы (isCancel=true).
// Связь с Sales по srid. При первом синке (нет курсора) — 90 дней, помечаем как notified.
func (s *Service) syncShopOrderRefusals(ctx context.Context, shop *domain.Shop) error {
	// WB отдаёт даты заказов/отказов без таймзоны в московском времени — парсим через глобальный wbReportLoc.
	dateFrom, err := s.orderSyncCursorRepo.Get(ctx, shop.ID)
	if err != nil {
		return fmt.Errorf("get order sync cursor: %w", err)
	}
	// Если wb_orders пустая (например, миграцию только что применили), нужно сделать backfill,
	// даже если курсор orders уже был установлен ранее.
	cnt, err := s.wbOrdersRepo.CountByShop(ctx, shop.ID)
	if err != nil {
		return fmt.Errorf("count wb_orders: %w", err)
	}
	isFirstSync := dateFrom.IsZero() || cnt == 0
	if isFirstSync {
		dateFrom = time.Now().AddDate(0, 0, -90)
	} else {
		// WB Statistics API задерживает данные на 1-4 часа: курсор сдвигаем назад,
		// чтобы не пропустить отказы, которые WB опубликовал с опозданием.
		// CreateBatch использует ON CONFLICT DO NOTHING — дублей не будет.
		const wbDataLagBuffer = 4 * time.Hour
		dateFrom = dateFrom.Add(-wbDataLagBuffer)
	}

	var notifiedAt *time.Time
	var refusalsTotal int
	if isFirstSync {
		now := time.Now()
		notifiedAt = &now
	}

	for {
		resp, err := s.wbClient.GetOrders(ctx, shop.APIToken, dateFrom, 0, nil)
		if err != nil {
			return fmt.Errorf("wb get orders: %w", err)
		}
		if len(resp) == 0 {
			break
		}
		// сохраняем все заказы для расчёта спроса (orders_net/day)
		now := time.Now()
		all := make([]*model.WBOrder, 0, len(resp))
		var batch []*domain.OrderRefusal
		var maxLastChange time.Time
		for i := range resp {
			item := &resp[i]
			dt, _ := item.ParseDate(nil)
			lc, _ := item.ParseLastChangeDate(nil)
			if lc.After(maxLastChange) {
				maxLastChange = lc
			}
			var cancelDt *time.Time
			if item.IsCancel {
				if t, err := item.ParseCancelDate(nil); err == nil && !t.IsZero() {
					cancelDt = &t
				}
			}
			if dt.IsZero() {
				dt = now
			}
			if lc.IsZero() {
				lc = now
			}
			all = append(all, &model.WBOrder{
				ShopID:          shop.ID,
				Srid:            item.Srid,
				NmID:            item.NmID,
				Barcode:         item.Barcode,
				SupplierArticle: item.SupplierArticle,
				Date:            dt,
				LastChangeDate:  lc,
				IsCancel:        item.IsCancel,
				CancelDate:      cancelDt,
				CreatedAt:       now,
				UpdatedAt:       now,
			})
			if !item.IsCancel {
				continue
			}
			r := wbmapper.OrderItemToRefusal(*item, shop.ID, nil)
			r.NotifiedAt = notifiedAt
			batch = append(batch, &r)
		}
		if err := s.wbOrdersRepo.UpsertBatch(ctx, all); err != nil {
			return fmt.Errorf("upsert wb_orders batch: %w", err)
		}
		if len(batch) > 0 {
			refusalsTotal += len(batch)
			if err := s.orderRefusalsRepo.CreateBatch(ctx, batch); err != nil {
				return fmt.Errorf("create order_refusals batch: %w", err)
			}
		}
		if !maxLastChange.IsZero() {
			dateFrom = maxLastChange.Add(time.Millisecond)
			if err := s.orderSyncCursorRepo.Set(ctx, shop.ID, maxLastChange); err != nil {
				return fmt.Errorf("set order sync cursor: %w", err)
			}
		} else if len(resp) > 0 {
			if lastChange, _ := resp[len(resp)-1].ParseLastChangeDate(nil); !lastChange.IsZero() {
				dateFrom = lastChange.Add(time.Millisecond)
				_ = s.orderSyncCursorRepo.Set(ctx, shop.ID, lastChange)
			}
		}
	}

	if isFirstSync {
		now := time.Now()
		if err := s.orderRefusalsRepo.MarkAllNotifiedByShop(ctx, shop.ID, now); err != nil {
			return fmt.Errorf("mark order_refusals notified: %w", err)
		}
		if refusalsTotal > 0 {
			msg := formatFirstSyncRefusalsSummaryRU(refusalsTotal)
			chatIDs, err := s.endpointsRepo.ListEnabledTelegramChatIDsByShopID(ctx, shop.ID)
			if err != nil {
				s.log.Warn("first sync refusals summary: list endpoints failed", "shop_id", shop.ID, "err", err)
			}
			if len(chatIDs) == 0 {
				user, _ := s.userRepo.GetByID(ctx, shop.UserID)
				if user != nil && user.TelegramChatID != nil && *user.TelegramChatID != 0 {
					chatIDs = []int64{*user.TelegramChatID}
				}
			}
			for _, chatID := range chatIDs {
				if err := s.tg.SendMessage(ctx, chatID, msg); err != nil {
					s.log.Warn("first sync refusals summary: send failed", "shop_id", shop.ID, "chat_id", chatID, "err", err)
				}
			}
		}
	}

	return nil
}

func formatFirstSyncRefusalsSummaryRU(refusalsCount int) string {
	return fmt.Sprintf(
		"✅ Синхронизация отказов завершена!\n\n"+
			"🚫 Загружено отказов за последние 90 дней: %d\n\n"+
			"Теперь вы будете получать уведомления об отказах (не забрали) в реальном времени.",
		refusalsCount,
	)
}

func formatFirstSyncSalesSummaryRU(buyouts, returns int, revenue float64) string {
	return fmt.Sprintf(
		"✅ Синхронизация завершена!\n\n"+
			"📊 Загружено за последние 90 дней:\n"+
			"- Выкупов: %d\n"+
			"- Возвратов: %d\n"+
			"- Выручка: %s\n\n"+
			"Теперь вы будете получать уведомления о новых продажах в реальном времени.",
		buyouts, returns, formatRevenueRU(revenue),
	)
}

func formatRevenueRU(v float64) string {
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}
	s := fmt.Sprintf("%.0f", v)
	if len(s) > 3 {
		var b []byte
		for i, c := range s {
			if i > 0 && (len(s)-i)%3 == 0 {
				b = append(b, ',')
			}
			b = append(b, byte(c))
		}
		s = string(b)
	}
	return sign + s + "₽"
}

// syncShopProductsIncremental синхронизирует карточки WB инкрементально, используя
// курсор из предыдущего синка. Раз в fullCardsSyncInterval делает полный синк.
// Возвращает новый курсор (nil если карточек не было), флаг isFullSync и ошибку.
func (s *Service) syncShopProductsIncremental(ctx context.Context, shop *domain.Shop) (*dto.ResponseCursor, bool, error) {
	isFullSync := shop.LastFullCardsSyncAt == nil || time.Since(*shop.LastFullCardsSyncAt) >= fullCardsSyncInterval

	var startCursor *dto.ResponseCursor
	if !isFullSync && shop.CardsCursorUpdatedAt != "" {
		startCursor = &dto.ResponseCursor{
			UpdatedAt: shop.CardsCursorUpdatedAt,
			NmID:      shop.CardsCursorNmID,
		}
	}

	req := dto.GetCardsRequest{
		Settings: dto.Settings{
			Cursor: dto.Cursor{Limit: 100},
			Filter: dto.Filter{WithPhoto: -1},
		},
	}

	cards, newCursor, err := s.wbClient.GetAllCards(ctx, shop.APIToken, req, startCursor)
	if err != nil {
		return nil, false, fmt.Errorf("wb get cards: %w", err)
	}

	for _, card := range cards {
		p := wbmapper.CardToProduct(card, shop.ID, wbProductBaseURL)

		if err := s.productsRepo.Upsert(ctx, &p); err != nil {
			return nil, false, fmt.Errorf("upsert product external_id=%s: %w", p.ExternalID, err)
		}

		for _, sz := range card.Sizes {
			chrtID := strconv.Itoa(sz.ChrtID)

			var tech *string
			if sz.TechSize != "" {
				v := sz.TechSize
				tech = &v
			}

			ps := domain.ProductSize{
				ProductID:    p.ID,
				ExternalID:   chrtID,
				TechSizeName: tech,
			}

			if err := s.productSizesRepo.Upsert(ctx, &ps); err != nil {
				return nil, false, fmt.Errorf("upsert product_size product_id=%d chrt_id=%s: %w", p.ID, chrtID, err)
			}

			if err := s.productSizeSKUsRepo.ReplaceByProductSizeID(ctx, ps.ID, sz.Skus); err != nil {
				return nil, false, fmt.Errorf("replace product_size_skus size_id=%d: %w", ps.ID, err)
			}
		}
	}

	return newCursor, isFullSync, nil
}

func (s *Service) syncShopSnapshots(
	ctx context.Context,
	shop *domain.Shop,
	products []domain.Product,
) error {
	// 1) nmID -> product, ids, productIDs
	ids := make([]int64, 0, len(products))
	prodByNm := make(map[int64]domain.Product, len(products))
	productIDs := make([]int, 0, len(products))

	for _, p := range products {
		nmID, err := strconv.ParseInt(p.ExternalID, 10, 64)
		if err != nil {
			s.log.Warn("invalid external_id, skip for snapshots",
				"product_id", p.ID, "external_id", p.ExternalID, "err", err)
			continue
		}
		ids = append(ids, nmID)
		prodByNm[nmID] = p
		productIDs = append(productIDs, p.ID)
	}
	if len(ids) == 0 {
		return nil
	}

	// 2) prices
	pricesResp, err := s.wbClient.GetPricesByNmIDs(ctx, shop.APIToken, ids)
	if err != nil {
		return fmt.Errorf("wb get prices: %w", err)
	}

	// 3) seller warehouses
	sellerWhs, err := s.sellerWarehousesRepo.ListByShop(ctx, shop.ID)
	if err != nil {
		return fmt.Errorf("list seller warehouses: %w", err)
	}

	// 4) sizes for products
	sizes, err := s.productSizesRepo.ListByProductIDs(ctx, productIDs)
	if err != nil {
		return fmt.Errorf("list product sizes: %w", err)
	}
	// sizesByProduct: product_id -> []ProductSize
	// sizeIDToChrtByProduct: product_id -> map[product_size_id]chrtID
	// product_id -> product_size_id -> chrtID(string)
	sizeIDToChrtByProduct := make(map[int]map[int]string)

	// уникальные chrtID для WB stocks
	chrtIDs := make([]int64, 0, len(sizes))
	seenChrt := make(map[int64]struct{})

	for _, sz := range sizes {
		if sz.ExternalID == "" {
			continue
		}

		if _, ok := sizeIDToChrtByProduct[sz.ProductID]; !ok {
			sizeIDToChrtByProduct[sz.ProductID] = make(map[int]string)
		}
		sizeIDToChrtByProduct[sz.ProductID][sz.ID] = sz.ExternalID

		// parse chrtID ТОЛЬКО ЗДЕСЬ
		chrtID, err := strconv.ParseInt(sz.ExternalID, 10, 64)
		if err != nil {
			s.log.Warn("invalid chrtID in product size, skip for stocks",
				"product_size_id", sz.ID,
				"external_id", sz.ExternalID,
				"err", err,
			)
			continue
		}

		if _, ok := seenChrt[chrtID]; ok {
			continue
		}
		seenChrt[chrtID] = struct{}{}
		chrtIDs = append(chrtIDs, chrtID)
	}

	// 5) stocks by warehouse
	stockByChrtID := make(map[int64]int)
	stocksUnavailable := false

	if len(sellerWhs) > 0 && len(chrtIDs) > 0 {
		for _, wh := range sellerWhs {
			items, err := s.wbClient.GetStocksByChrtIDs(
				ctx,
				shop.APIToken,
				int64(wh.WBID),
				chrtIDs,
			)
			if err != nil {
				s.log.Error("get stocks failed",
					"shop_id", shop.ID,
					"warehouse_id", wh.ID,
					"chrtIDs_count", len(chrtIDs),
					"err", err,
				)
				stocksUnavailable = true
				break
			}
			for _, st := range items {
				stockByChrtID[st.ChrtID] += st.Amount
			}
		}
	}

	// 8) stock total by product + by nm
	var stockByNm map[int64]int
	if !stocksUnavailable {
		stockByProductID := make(map[int]int)

		for _, sz := range sizes {
			chrtID, err := strconv.ParseInt(sz.ExternalID, 10, 64)
			if err != nil {
				continue
			}
			stockByProductID[sz.ProductID] += stockByChrtID[chrtID]
		}

		stockByNm = make(map[int64]int, len(prodByNm))
		for nmID, p := range prodByNm {
			stockByNm[nmID] = stockByProductID[p.ID]
		}
	}

	// 9) product snapshots (only product-level)
	items := wbmapper.ToSnapshots(prodByNm, pricesResp, stockByNm)

	// 10) Build prices map: product_id + chrtID -> wbSizePrice
	type wbSizePrice struct {
		BasePrice    *float64
		FinalPrice   *float64
		ClubPrice    *float64
		Discount     *float64
		ClubDiscount *float64
		Extra        []byte
	}
	sizeByProductAndTech := make(map[int]map[string]int)

	for _, sz := range sizes {
		if sz.TechSizeName == nil {
			continue
		}
		if _, ok := sizeByProductAndTech[sz.ProductID]; !ok {
			sizeByProductAndTech[sz.ProductID] = make(map[string]int)
		}
		sizeByProductAndTech[sz.ProductID][*sz.TechSizeName] = sz.ID
	}

	pricesByProductAndSize := make(map[int]map[int]wbSizePrice)

	for _, g := range pricesResp {
		p, ok := prodByNm[g.NmID]
		if !ok {
			continue
		}
		// if _, ok := pricesByProductAndChrt[p.ID]; !ok {
		// 	pricesByProductAndChrt[p.ID] = make(map[string]wbSizePrice)
		// }

		sizeIndex := sizeByProductAndTech[p.ID]
		if len(sizeIndex) == 0 {
			continue
		}

		if _, ok := pricesByProductAndSize[p.ID]; !ok {
			pricesByProductAndSize[p.ID] = make(map[int]wbSizePrice)
		}

		disc := g.Discount
		clubDisc := g.ClubDiscount

		for _, sz := range g.Sizes {
			psID, ok := sizeIndex[sz.TechSizeName]
			if !ok {
				// размер есть в prices, но нет в cards
				continue
			}
			base := sz.Price
			final := sz.DiscountedPrice
			club := sz.ClubDiscountedPrice

			pricesByProductAndSize[p.ID][psID] = wbSizePrice{
				BasePrice:    &base,
				FinalPrice:   &final,
				ClubPrice:    &club,
				Discount:     &disc,
				ClubDiscount: &clubDisc,
			}
		}
	}

	for _, item := range items {
		snap := item.Snapshot

		if err := s.snapshotsRepo.Create(ctx, &snap); err != nil {
			return fmt.Errorf("create snapshot for product %d: %w", snap.ProductID, err)
		}

		// 11.1)  prices per size

		//priceMap := pricesByProductAndChrt[snap.ProductID] // chrtID -> price
		priceMap := pricesByProductAndSize[snap.ProductID]

		for psID, p := range priceMap {
			priceSnap := domain.ProductSizePriceSnapshot{
				ProductSnapshotID: snap.ID,
				ProductSizeID:     psID,
				BasePrice:         p.BasePrice,
				FinalPrice:        p.FinalPrice,
				ClubPrice:         p.ClubPrice,
				Discount:          p.Discount,
				ClubDiscount:      p.ClubDiscount,
				Extra:             p.Extra,
			}

			if err := s.sizePriceSnapshotsRepo.Create(ctx, &priceSnap); err != nil {
				return fmt.Errorf("create size price snapshot (product %d, size %d): %w", snap.ProductID, psID, err)
			}
		}

		// 11.2) stocks per size (iterate sizes, not skuRows!)
		if !stocksUnavailable {
			sizeMap := sizeIDToChrtByProduct[snap.ProductID]
			for psID, chrtIDStr := range sizeMap {
				chrtID, err := strconv.ParseInt(chrtIDStr, 10, 64)
				if err != nil {
					continue
				}

				q := stockByChrtID[chrtID]
				qq := q

				sizeSnap := domain.ProductSizeSnapshot{
					ProductSnapshotID: snap.ID,
					ProductSizeID:     psID,
					Stock:             &qq,
				}

				if err := s.sizeSnapshotsRepo.Create(ctx, &sizeSnap); err != nil {
					return fmt.Errorf("create size stock snapshot (product %d, size %d): %w",
						snap.ProductID, psID, err)
				}
			}
		}
	}

	return nil
}
