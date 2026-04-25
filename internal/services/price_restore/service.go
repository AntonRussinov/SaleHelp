package price_restore

import (
	"context"
	"log/slog"
	"time"

	"salehelp/domain"
	wbclient "salehelp/internal/adapter/wb/client"
	"salehelp/pkg/logger/sl"
)

const maxGoodsPerRequest = 1000

type (
	BatchRepo interface {
		GetByID(ctx context.Context, id int) (*domain.PriceRestoreBatch, error)
		GetTriggeredIDsByBatchID(ctx context.Context, batchID int) ([]int, error)
		SetRestoredAt(ctx context.Context, batchID int, at time.Time) error
		TrySetRestoring(ctx context.Context, batchID int) (bool, error)
		SetRestoreFailed(ctx context.Context, batchID int) error
	}
	TriggeredRepo interface {
		GetByIDs(ctx context.Context, ids []int) (map[int]*domain.AlertTriggered, error)
	}
	ProductSnapshotRepo interface {
		GetByID(ctx context.Context, id int) (*domain.ProductSnapshot, error)
		GetPreviousBySnapshotID(ctx context.Context, snapshotID int) (*domain.ProductSnapshot, error)
	}
	ProductSizePriceSnapshotRepo interface {
		ListBySnapshotID(ctx context.Context, snapshotID int) ([]*domain.ProductSizePriceSnapshot, error)
	}
	ProductRepo interface {
		GetByIDs(ctx context.Context, ids []int) (map[int]*domain.Product, error)
	}
	ProductSizeRepo interface {
		ListByProductIDs(ctx context.Context, productIDs []int) ([]*domain.ProductSize, error)
	}
	ShopRepo interface {
		GetByID(ctx context.Context, id int) (*domain.Shop, error)
	}
)

type Service struct {
	log          *slog.Logger
	batchRepo    BatchRepo
	triggered    TriggeredRepo
	snapRepo     ProductSnapshotRepo
	sizeSnapRepo ProductSizePriceSnapshotRepo
	productRepo  ProductRepo
	sizeRepo     ProductSizeRepo
	shopRepo     ShopRepo
	wbClient     *wbclient.Client
}

func NewService(
	log *slog.Logger,
	batchRepo BatchRepo,
	triggered TriggeredRepo,
	snapRepo ProductSnapshotRepo,
	sizeSnapRepo ProductSizePriceSnapshotRepo,
	productRepo ProductRepo,
	sizeRepo ProductSizeRepo,
	shopRepo ShopRepo,
	wbClient *wbclient.Client,
) *Service {
	return &Service{
		log:          log.With(sl.Module("PriceRestore")),
		batchRepo:    batchRepo,
		triggered:    triggered,
		snapRepo:     snapRepo,
		sizeSnapRepo: sizeSnapRepo,
		productRepo:  productRepo,
		sizeRepo:     sizeRepo,
		shopRepo:     shopRepo,
		wbClient:     wbClient,
	}
}

// RestoreResult — результат попытки восстановления.
type RestoreResult struct {
	ShopID          int
	RestoredCount   int
	AlreadyRestored bool
	AlreadyRestoring bool
	RestoredAt      *time.Time
	Error           error
}
