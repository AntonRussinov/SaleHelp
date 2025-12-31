package monitoring

import (
	"context"
	"log/slog"
	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
	"salehelp/internal/storage/pggorm/repository/productsizesku"
)

type (
	Service struct {
		shopsRepo              ShopRepository
		productsRepo           ProductRepository
		snapshotsRepo          ProductSnapshotRepository
		sizeSnapshotsRepo      ProductSizeSnapshotRepository
		sizePriceSnapshotsRepo ProductSizePriceRepository
		productSizesRepo       ProductSizeRepository
		productSizeSKUsRepo    ProductSizeSKURepository
		wbWarehousesRepo       WBWarehousesRepository
		sellerWarehousesRepo   SellerWarehousesRepository

		wbClient WBClient
		log      *slog.Logger
	}

	// Локальный интерфейс под нужды сервиса
	ShopRepository interface {
		GetByID(ctx context.Context, id int) (*domain.Shop, error)
		ListAll(ctx context.Context) ([]*domain.Shop, error)
	}

	ProductRepository interface {
		Upsert(ctx context.Context, p *domain.Product) error
		ListIDsByShopID(ctx context.Context, shopID int) ([]int, error)
	}

	ProductSnapshotRepository interface {
		Create(ctx context.Context, s *domain.ProductSnapshot) error
		GetLatestByProductID(ctx context.Context, productID int) (*domain.ProductSnapshot, error)
		GetLatestByProductIDs(ctx context.Context, productIDs []int) (map[int]*domain.ProductSnapshot, error)
	}

	ProductSizeSnapshotRepository interface {
		Create(ctx context.Context, s *domain.ProductSizeSnapshot) error
		ListBySnapshotID(ctx context.Context, snapshotID int) ([]*domain.ProductSizeSnapshot, error)
		ListBySnapshotIDs(ctx context.Context, snapshotIDs []int) ([]*domain.ProductSizeSnapshot, error)
	}

	ProductSizeRepository interface {
		Upsert(ctx context.Context, ps *domain.ProductSize) error
		FindByProductAndExternalID(ctx context.Context, productID int, externalID string) (*domain.ProductSize, error)
		ListByProductID(ctx context.Context, productID int) ([]*domain.ProductSize, error)
		ListByProductIDs(ctx context.Context, productIDs []int64) ([]*domain.ProductSize, error)
	}

	ProductSizePriceRepository interface {
		Create(ctx context.Context, s *domain.ProductSizePriceSnapshot) error
		ListBySnapshotID(ctx context.Context, snapshotID int) ([]*domain.ProductSizePriceSnapshot, error)
		ListBySnapshotIDs(ctx context.Context, snapshotIDs []int) ([]*domain.ProductSizePriceSnapshot, error)
	}

	ProductSizeSKURepository interface {
		ReplaceByProductSizeID(ctx context.Context, productSizeID int, skus []string) error
		ListByProductSizeIDs(ctx context.Context, productSizeIDs []int) (map[int][]string, error)
		ListByProductIDs(ctx context.Context, productIDs []int64) ([]productsizesku.SKUWithRefs, error)
	}

	WBWarehousesRepository interface {
		Upsert(ctx context.Context, w *domain.WBWarehouse) error
		// при желании ещё List, GetByWBID и т.п.
	}

	SellerWarehousesRepository interface {
		ReplaceForShop(ctx context.Context, shopID int64, warehouses []domain.SellerWarehouse) error
		ListByShop(ctx context.Context, shopID int) ([]*domain.SellerWarehouse, error)
	}

	// Локальный интерфейс под WB-клиент
	WBClient interface {
		GetCards(ctx context.Context, apiKey string, req dto.GetCardsRequest) (*dto.GetCardsResponse, error)
		GetAllCards(ctx context.Context, apiKey string, baseReq dto.GetCardsRequest) ([]dto.Card, error)
		GetPrices(ctx context.Context, token string, req dto.GetPricesRequest) (*dto.GetPricesResponse, error)
		GetPricesByNmIDs(ctx context.Context, apiKey string, nmIDs []int64) ([]dto.PriceGood, error)
		GetStocks(ctx context.Context, apiKey string, warehouseID int64, req dto.GetStocksRequest) (*dto.GetStocksResponse, error)
		GetStocksByChrtIDs(ctx context.Context, apiKey string, warehouseID int64, chrtIDs []int64) ([]dto.StockItem, error)
		GetWBWarehouses(ctx context.Context, apiKey string) ([]dto.WBWarehouse, error)
		GetSellerWarehouses(ctx context.Context, apiKey string) ([]dto.SellerWarehouse, error)
	}
)

func NewMonitoringService(
	shops ShopRepository,
	products ProductRepository,
	snapshots ProductSnapshotRepository,
	sizesnapshots ProductSizeSnapshotRepository,
	sizePriceSnapshots ProductSizePriceRepository,
	productSizes ProductSizeRepository,
	productSizeSKUs ProductSizeSKURepository,
	wbWarehousesRepo WBWarehousesRepository,
	sellerWarehousesRepo SellerWarehousesRepository,
	wb WBClient,
	log *slog.Logger,
) *Service {
	return &Service{
		shopsRepo:              shops,
		productsRepo:           products,
		snapshotsRepo:          snapshots,
		sizeSnapshotsRepo:      sizesnapshots,
		sizePriceSnapshotsRepo: sizePriceSnapshots,
		productSizesRepo:       productSizes,
		productSizeSKUsRepo:    productSizeSKUs,
		wbWarehousesRepo:       wbWarehousesRepo,
		sellerWarehousesRepo:   sellerWarehousesRepo,
		wbClient:               wb,
		log:                    log,
	}
}
