package alert

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/services/subscription"
)

const (
	AlertTypeStockLow  = "stock_low"
	AlertTypeStockZero = "stock_zero"
	AlertTypePriceDrop = "price_drop"
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

func NewAlertEngine(
	subSvc SubscriptionService,
	alertsRepo AlertsRepo,
	alertsTriggeredRepo AlertsTriggeredRepo,
	productsRepo ProductsRepo,
	productSizesRepo ProductSizesRepo,
	productSizeSnapshotsRepo ProductSizeSnapshotsRepo,
	snapshotsRepo ProductSnapshotsRepo,
	deliveriesFanoutRepo DeliveriesFanoutRepo,
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
	}
}
