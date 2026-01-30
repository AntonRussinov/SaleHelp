package alert

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"salehelp/domain"
	"salehelp/internal/services/notification"
	"salehelp/internal/services/subscription"
)

type stockLowCfg struct {
	Threshold int `json:"threshold"`
}

type AlertEngine struct {
	subSvc                   SubscriptionService
	alertsRepo               AlertsRepo
	alertsTriggeredRepo      AlertsTriggeredRepo
	productsRepo             ProductsRepo
	productSizesRepo         ProductSizesRepo
	productSizeSnapshotsRepo ProductSizeSnapshotsRepo
	snapshotsRepo            ProductSnapshotsRepo
	deliveriesFanoutRepo     DeliveriesFanoutRepo
	fbsOrdersRepo            FBSOrderRepo
	endpointsRepo            EndpointsRepo
	tg                       TelegramSender
	log                      *slog.Logger
}

type SubscriptionService interface {
	GetShopAccessLevel(ctx context.Context, shopID int) (subscription.AccessLevel, *domain.Subscription, error)
}

type AlertsRepo interface {
	//ListActiveByShop(ctx context.Context, shopID int) ([]domain.Alert, error)
	ListActiveByShop(ctx context.Context, shopID int) ([]*domain.Alert, error)
}

type AlertsTriggeredRepo interface {
	//ExistsForSnapshot(ctx context.Context, alertID, productID, snapshotID int) (bool, error)
	CreateIfNotExists(ctx context.Context, at *domain.AlertTriggered) (id int, inserted bool, err error)

	//Create(ctx context.Context, at *domain.AlertTriggered) error
}

type ProductSnapshotsRepo interface {
	GetLatestByProductIDsLimit(ctx context.Context, productIDs []int, limit int) (map[int][]*domain.ProductSnapshot, error)
}

type ProductSizesRepo interface {
	ListByProductIDs(ctx context.Context, productIDs []int) ([]*domain.ProductSize, error)
}

type ProductSizeSnapshotsRepo interface {
	// последние N снэпшотов по каждому product_size_id
	GetLatestBySizeIDsLimitMap(ctx context.Context, sizeIDs []int, limit int) (map[int][]*domain.ProductSizeSnapshot, error)
}

type ProductsRepo interface {
	ListByShop(ctx context.Context, shopID int, limit, offset int) ([]*domain.Product, error)
}

type DeliveriesFanoutRepo interface {
	CreateDeliveriesForTriggered(ctx context.Context, triggeredID int, shopID int) error
}

type FBSOrderRepo interface {
	ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.FBSOrder, error)
	MarkNotified(ctx context.Context, orderIDs []int, at time.Time) error
}

type EndpointsRepo interface {
	ListByShopID(ctx context.Context, shopID int, onlyEnabled bool) ([]*domain.NotificationEndpoint, error)
}

type TelegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

func NewAlertEngine(
	subSvc SubscriptionService,
	alertsRepo AlertsRepo,
	alertsTriggeredRepo AlertsTriggeredRepo,
	productsRepo ProductsRepo,
	productSizesRepo ProductSizesRepo,
	productSizeSnapshotsRepo ProductSizeSnapshotsRepo,
	snapshotsRepo ProductSnapshotsRepo,
	deliveriesFanoutRepo DeliveriesFanoutRepo,
	fbsOrdersRepo FBSOrderRepo,
	endpointsRepo EndpointsRepo,
	tg TelegramSender,
	log *slog.Logger,
) *AlertEngine {
	return &AlertEngine{
		subSvc:                   subSvc,
		alertsRepo:               alertsRepo,
		alertsTriggeredRepo:      alertsTriggeredRepo,
		productsRepo:             productsRepo,
		productSizesRepo:         productSizesRepo,
		productSizeSnapshotsRepo: productSizeSnapshotsRepo,
		snapshotsRepo:            snapshotsRepo,
		deliveriesFanoutRepo:     deliveriesFanoutRepo,
		fbsOrdersRepo:            fbsOrdersRepo,
		endpointsRepo:            endpointsRepo,
		tg:                       tg,
		log:                      log,
	}
}

// runShopFBSOrders отправляет уведомления о новых FBS-заказах в Telegram и помечает заказы как уведомлённые.
// products — уже загруженные товары магазина (из RunShop); если nil/пусто и есть заказы — подгружаем ListByShop для подстановки названий.
func (e *AlertEngine) runShopFBSOrders(ctx context.Context, shopID int, products []*domain.Product) error {
	orders, err := e.fbsOrdersRepo.ListUnnotifiedByShop(ctx, shopID, 50)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		return nil
	}
	// Названия товаров: используем переданные products или один раз подгружаем ListByShop
	if len(products) == 0 {
		products, err = e.productsRepo.ListByShop(ctx, shopID, 0, 0)
		if err != nil {
			return err
		}
	}
	nmIDSet := make(map[int64]struct{}, len(orders))
	for _, o := range orders {
		nmIDSet[o.NmID] = struct{}{}
	}
	titleByNmID := make(map[int64]string)
	for _, p := range products {
		if p == nil || p.Title == "" || p.ExternalID == "" {
			continue
		}
		nmID, err := strconv.ParseInt(p.ExternalID, 10, 64)
		if err != nil {
			continue
		}
		if _, ok := nmIDSet[nmID]; ok {
			titleByNmID[nmID] = p.Title
		}
	}
	endpoints, err := e.endpointsRepo.ListByShopID(ctx, shopID, true)
	if err != nil {
		return err
	}
	text := notifications.FormatFBSOrdersMessageRU(orders, titleByNmID)
	now := time.Now()
	sent := false
	for _, ep := range endpoints {
		if ep.Channel != "telegram" || !ep.Enabled {
			continue
		}
		chatID, err := strconv.ParseInt(ep.Address, 10, 64)
		if err != nil {
			e.log.Warn("fbs orders: invalid telegram chat_id", "endpoint_id", ep.ID, "address", ep.Address, "err", err)
			continue
		}
		if err := e.tg.SendMessage(ctx, chatID, text); err != nil {
			e.log.Error("fbs orders: send failed", "endpoint_id", ep.ID, "err", err)
			continue
		}
		sent = true
	}
	if sent {
		ids := make([]int, 0, len(orders))
		for _, o := range orders {
			ids = append(ids, o.ID)
		}
		if err := e.fbsOrdersRepo.MarkNotified(ctx, ids, now); err != nil {
			return err
		}
	}
	return nil
}
