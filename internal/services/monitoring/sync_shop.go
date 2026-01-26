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

func (s *Service) SyncAllShops(ctx context.Context) error {
	shops, err := s.shopsRepo.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("sync all shops: list shops: %w", err)
	}

	for _, shop := range shops {
		if err := s.SyncShop(ctx, shop.ID); err != nil {
			// логируем, но не фейлим весь процесс
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

	// единый апдейт статуса синка при любом выходе
	defer func() {
		now := time.Now()

		// всегда фиксируем попытку
		upd := model.SyncStatusUpdate{
			LastSyncAt: &now,
		}

		if err == nil {
			upd.LastSyncOkAt = &now
			upd.LastSyncStatusCode = nil
			upd.LastSyncError = nil
			upd.LastSyncFailedAt = nil
		} else {
			upd.LastSyncFailedAt = &now

			if sc, ok := wbStatusCode(err); ok {
				upd.LastSyncStatusCode = &sc
			}
			m := err.Error()
			upd.LastSyncError = &m
		}

		_ = s.shopsRepo.UpdateSyncStatus(ctx, shopID, upd)
	}()

	if err := s.syncShopWarehouses(ctx, shop); err != nil {
		return fmt.Errorf("sync shop: sync warehouses: %w", err)
	}

	// 1) карточки -> products
	products, err := s.syncShopProducts(ctx, shop)
	if err != nil {
		return fmt.Errorf("sync shop: sync products: %w", err)
	}
	if len(products) == 0 {
		// нет товаров – нет снапшотов
		return nil
	}

	// 2) цены/остатки -> snapshots
	if err = s.syncShopSnapshots(ctx, shop, products); err != nil {
		return fmt.Errorf("sync shop: sync snapshots: %w", err)
	}

	// aha-moment (первый синк) — оставляем как у тебя
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

// syncShopProducts синхронизирует карточки WB в таблицу products
// и возвращает актуальный список продуктов этого магазина.
func (s *Service) syncShopProducts(ctx context.Context, shop *domain.Shop) ([]domain.Product, error) {
	req := dto.GetCardsRequest{
		Settings: dto.Settings{
			Cursor: dto.Cursor{
				Limit: 100,
			},
			Filter: dto.Filter{
				TextSearch:            "",
				AllowedCategoriesOnly: false,
				TagIDs:                []int{},
				ObjectIDs:             []int{},
				Brands:                []string{},
				ImtID:                 0,
				WithPhoto:             -1,
			},
		},
	}

	cards, err := s.wbClient.GetAllCards(ctx, shop.APIToken, req)
	if err != nil {
		return nil, fmt.Errorf("wb get cards: %w", err)
	}

	if len(cards) == 0 {
		return nil, nil
	}

	products := make([]domain.Product, 0, len(cards))
	for _, card := range cards {
		p := wbmapper.CardToProduct(card, shop.ID, wbProductBaseURL)

		if err := s.productsRepo.Upsert(ctx, &p); err != nil {
			return nil, fmt.Errorf("upsert product external_id=%s: %w",
				p.ExternalID, err)
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
				return nil, fmt.Errorf("upsert product_size product_id=%d chrt_id=%s: %w", p.ID, chrtID, err)
			}

			// Replace — важно, чтобы старые SKUs удалялись
			if err := s.productSizeSKUsRepo.ReplaceByProductSizeID(ctx, ps.ID, sz.Skus); err != nil {
				return nil, fmt.Errorf("replace product_size_skus size_id=%d: %w", ps.ID, err)
			}
		}

		// важно: Upsert через GORM обычно проставляет p.ID
		products = append(products, p)
	}

	return products, nil
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

	// 4) sku rows (product_id + product_size_id + sku)
	// skuRows, err := s.productSizeSKUsRepo.ListByProductIDs(ctx, productIDs)
	// if err != nil {
	// 	return fmt.Errorf("list product size skus: %w", err)
	// }

	// 5) sizes for products (product_size_id -> chrtID) сгруппировать по product_id
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

	// 7) stocks by sku
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
				// Важно: если WB stocks недоступен, нельзя трактовать это как "0" и перезаписывать сток.
				// Иначе получим ложные алерты "всё закончилось".
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
