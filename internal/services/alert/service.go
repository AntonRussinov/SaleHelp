package alert

import (
	"context"
	"log/slog"
	"time"

	"salehelp/domain"
	"salehelp/internal/services/notification_settings"
	"salehelp/internal/services/subscription"
)

type stockLowCfg struct {
	Threshold int `json:"threshold"`
}

type AlertEngine struct {
	subSvc                   SubscriptionService
	notifSettings            *notification_settings.Service
	alertsRepo               AlertsRepo
	alertsTriggeredRepo      AlertsTriggeredRepo
	productsRepo             ProductsRepo
	productSizesRepo         ProductSizesRepo
	productSizeSnapshotsRepo ProductSizeSnapshotsRepo
	snapshotsRepo            ProductSnapshotsRepo
	deliveriesFanoutRepo     DeliveriesFanoutRepo
	fbsOrdersRepo            FBSOrderRepo
	dbwOrdersRepo            DBWOrderRepo
	dbsOrdersRepo            DBSOrderRepo
	saleEventsRepo           SaleEventRepo
	orderRefusalsRepo        OrderRefusalRepo
	wbFeedbacksRepo          WBFeedbackRepo
	wbQuestionsRepo          WBQuestionRepo
	endpointsRepo            EndpointsRepo
	tg                       TelegramSender
	abcSnapshotRepo          ABCSnapshotRepo
	shopTimezoneRepo         ShopTimezoneRepo
	defaultLoc               *time.Location
	log                      *slog.Logger
}

// shopLocation возвращает timezone магазина из таблицы shops.
// Если у магазина не задан WBReportTimezone — использует defaultLoc (пояс из конфига, тот же что при парсинге дат WB).
func (e *AlertEngine) shopLocation(ctx context.Context, shopID int) *time.Location {
	fallback := e.defaultLoc
	if fallback == nil {
		fallback = time.UTC
	}
	if e.shopTimezoneRepo == nil {
		return fallback
	}
	shop, err := e.shopTimezoneRepo.GetByID(ctx, shopID)
	if err != nil || shop == nil || shop.WBReportTimezone == "" {
		return fallback
	}
	loc, err := time.LoadLocation(shop.WBReportTimezone)
	if err != nil {
		return fallback
	}
	return loc
}

func (e *AlertEngine) notifEnabled(ctx context.Context, shopID int, key string) bool {
	if e == nil || e.notifSettings == nil {
		return true
	}
	enabled, err := e.notifSettings.EffectiveEnabled(ctx, shopID, key)
	if err != nil {
		if e.log != nil {
			e.log.Warn("notification settings check failed, treat as enabled", "shop_id", shopID, "key", key, "err", err)
		}
		return true
	}
	return enabled
}

type WBFeedbackRepo interface {
	ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.WBFeedback, error)
	MarkNotified(ctx context.Context, ids []int, at time.Time) error
}

type WBQuestionRepo interface {
	ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.WBQuestion, error)
	MarkNotified(ctx context.Context, ids []int, at time.Time) error
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

type DBWOrderRepo interface {
	ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.DBWOrder, error)
	MarkNotified(ctx context.Context, orderIDs []int, at time.Time) error
}

type DBSOrderRepo interface {
	ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.DBSOrder, error)
	MarkNotified(ctx context.Context, orderIDs []int, at time.Time) error
}

type SaleEventRepo interface {
	ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.SaleEvent, error)
	MarkNotified(ctx context.Context, eventIDs []int, at time.Time) error
}

type OrderRefusalRepo interface {
	ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.OrderRefusal, error)
	MarkNotified(ctx context.Context, ids []int, at time.Time) error
}

type EndpointsRepo interface {
	ListByShopID(ctx context.Context, shopID int, onlyEnabled bool) ([]*domain.NotificationEndpoint, error)
}

type TelegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

func NewAlertEngine(
	subSvc SubscriptionService,
	notifSettings *notification_settings.Service,
	alertsRepo AlertsRepo,
	alertsTriggeredRepo AlertsTriggeredRepo,
	productsRepo ProductsRepo,
	productSizesRepo ProductSizesRepo,
	productSizeSnapshotsRepo ProductSizeSnapshotsRepo,
	snapshotsRepo ProductSnapshotsRepo,
	deliveriesFanoutRepo DeliveriesFanoutRepo,
	fbsOrdersRepo FBSOrderRepo,
	dbwOrdersRepo DBWOrderRepo,
	dbsOrdersRepo DBSOrderRepo,
	saleEventsRepo SaleEventRepo,
	orderRefusalsRepo OrderRefusalRepo,
	wbFeedbacksRepo WBFeedbackRepo,
	wbQuestionsRepo WBQuestionRepo,
	endpointsRepo EndpointsRepo,
	tg TelegramSender,
	abcSnapshotRepo ABCSnapshotRepo,
	shopTimezoneRepo ShopTimezoneRepo,
	defaultLoc *time.Location,
	log *slog.Logger,
) *AlertEngine {
	return &AlertEngine{
		subSvc:                   subSvc,
		notifSettings:            notifSettings,
		alertsRepo:               alertsRepo,
		alertsTriggeredRepo:      alertsTriggeredRepo,
		productsRepo:             productsRepo,
		productSizesRepo:         productSizesRepo,
		productSizeSnapshotsRepo: productSizeSnapshotsRepo,
		snapshotsRepo:            snapshotsRepo,
		deliveriesFanoutRepo:     deliveriesFanoutRepo,
		fbsOrdersRepo:            fbsOrdersRepo,
		dbwOrdersRepo:            dbwOrdersRepo,
		dbsOrdersRepo:            dbsOrdersRepo,
		saleEventsRepo:           saleEventsRepo,
		orderRefusalsRepo:        orderRefusalsRepo,
		wbFeedbacksRepo:          wbFeedbacksRepo,
		wbQuestionsRepo:          wbQuestionsRepo,
		endpointsRepo:            endpointsRepo,
		tg:                       tg,
		abcSnapshotRepo:          abcSnapshotRepo,
		shopTimezoneRepo:         shopTimezoneRepo,
		defaultLoc:               defaultLoc,
		log:                      log,
	}
}

// ABCSnapshotRepo опционально: для приоритизации A-товаров в алертах по остаткам.
type ABCSnapshotRepo interface {
	GetLatestByShopAndPeriod(ctx context.Context, shopID, periodDays int) (*domain.ABCSnapshot, error)
}

type ShopTimezoneRepo interface {
	GetByID(ctx context.Context, id int) (*domain.Shop, error)
}
