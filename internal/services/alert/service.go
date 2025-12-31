package alert

import (
	"context"
	"salehelp/domain"
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
	alertsRepo           AlertsRepo
	alertsTriggeredRepo  AlertsTriggeredRepo
	productsRepo         ProductsRepo
	snapshotsRepo        ProductSnapshotsRepo
	deliveriesFanoutRepo DeliveriesFanoutRepo
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

type ProductsRepo interface {
	ListByShop(ctx context.Context, shopID int, limit, offset int) ([]*domain.Product, error)
}

type DeliveriesFanoutRepo interface {
	CreateDeliveriesForTriggered(ctx context.Context, triggeredID int, shopID int) error
}

func NewAlertEngine(
	alertsRepo AlertsRepo,
	alertsTriggeredRepo AlertsTriggeredRepo,
	productsRepo ProductsRepo,
	snapshotsRepo ProductSnapshotsRepo,
	deliveriesFanoutRepo DeliveriesFanoutRepo,
) *AlertEngine {
	return &AlertEngine{
		alertsRepo:           alertsRepo,
		alertsTriggeredRepo:  alertsTriggeredRepo,
		productsRepo:         productsRepo,
		snapshotsRepo:        snapshotsRepo,
		deliveriesFanoutRepo: deliveriesFanoutRepo,
	}
}
