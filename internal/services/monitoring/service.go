package monitoring

import (
	"context"
	"errors"
	"log/slog"
	"salehelp/domain"
	wbclient "salehelp/internal/adapter/wb/client"
	"salehelp/internal/adapter/wb/dto"
	"salehelp/internal/storage/pggorm/model"
	"salehelp/internal/storage/pggorm/repository/productsizesku"
	"time"
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
		fbsOrdersRepo          FBSOrderRepository
		dbwOrdersRepo          DBWOrderRepository
		dbsOrdersRepo          DBSOrderRepository
		saleEventsRepo         SaleEventRepository
		orderRefusalsRepo      OrderRefusalRepository
		wbOrdersRepo           WBOrdersRepository
		orderSyncCursorRepo    OrderSyncCursorRepository
		wbFeedbacksRepo        WBFeedbackRepository
		wbQuestionsRepo        WBQuestionRepository
		endpointsRepo          EndpointsRepo
		userRepo               UserRepo

		tg                TelegramAPI
		wbClient          WBClient
		defaultReportLoc  *time.Location // пояс по умолчанию для отчётов WB (из конфига), если у магазина не задан wb_report_timezone
		log               *slog.Logger
	}

	WBFeedbackRepository interface {
		Create(ctx context.Context, f *domain.WBFeedback) error
	}

	WBQuestionRepository interface {
		Create(ctx context.Context, q *domain.WBQuestion) error
	}

	// Локальный интерфейс под нужды сервиса
	ShopRepository interface {
		GetByID(ctx context.Context, id int) (*domain.Shop, error)
		ListAll(ctx context.Context) ([]*domain.Shop, error)
		UpdateLastSyncAt(ctx context.Context, shopID int, t time.Time) error
		UpdateSyncStatus(ctx context.Context, shopID int, update model.SyncStatusUpdate) error
		UpdateFirstSyncNotifiedAt(ctx context.Context, shopID int, t time.Time) error
	}

	ProductRepository interface {
		Upsert(ctx context.Context, p *domain.Product) error
		ListIDsByShopID(ctx context.Context, shopID int) ([]int, error)
		CountByShop(ctx context.Context, shopID int) (int, error)
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
		ListByProductIDs(ctx context.Context, productIDs []int) ([]*domain.ProductSize, error)
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

	FBSOrderRepository interface {
		Create(ctx context.Context, o *domain.FBSOrder) error
		ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.FBSOrder, error)
		MarkNotified(ctx context.Context, orderIDs []int, at time.Time) error
	}

	DBWOrderRepository interface {
		Create(ctx context.Context, o *domain.DBWOrder) error
		ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.DBWOrder, error)
		MarkNotified(ctx context.Context, orderIDs []int, at time.Time) error
	}

	DBSOrderRepository interface {
		Create(ctx context.Context, o *domain.DBSOrder) error
		ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.DBSOrder, error)
		MarkNotified(ctx context.Context, orderIDs []int, at time.Time) error
	}

	SaleEventRepository interface {
		Create(ctx context.Context, e *domain.SaleEvent) error
		CreateBatch(ctx context.Context, events []*domain.SaleEvent) error
		ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.SaleEvent, error)
		MarkNotified(ctx context.Context, eventIDs []int, at time.Time) error
		MarkAllNotifiedByShop(ctx context.Context, shopID int, at time.Time) error
		GetMaxLastChangeDateByShop(ctx context.Context, shopID int) (time.Time, error)
	}

	OrderRefusalRepository interface {
		CreateBatch(ctx context.Context, refusals []*domain.OrderRefusal) error
		MarkAllNotifiedByShop(ctx context.Context, shopID int, at time.Time) error
	}

	OrderSyncCursorRepository interface {
		Get(ctx context.Context, shopID int) (time.Time, error)
		Set(ctx context.Context, shopID int, lastChangeDate time.Time) error
	}

	WBOrdersRepository interface {
		UpsertBatch(ctx context.Context, rows []*model.WBOrder) error
		CountByShop(ctx context.Context, shopID int) (int, error)
	}

	UserRepo interface {
		GetByID(ctx context.Context, id int) (*domain.User, error)
	}

	EndpointsRepo interface {
		ListEnabledTelegramChatIDsByShopID(ctx context.Context, shopID int) ([]int64, error)
	}

	// Локальный интерфейс под WB-клиент
	WBClient interface {
		GetCards(ctx context.Context, apiKey string, req dto.GetCardsRequest) (*dto.GetCardsResponse, error)
		GetAllCards(ctx context.Context, apiKey string, baseReq dto.GetCardsRequest) ([]dto.Card, error)
		GetPrices(ctx context.Context, token string, req dto.GetPricesRequest) (*dto.GetPricesResponse, error)
		GetPricesByNmIDs(ctx context.Context, apiKey string, nmIDs []int64) ([]dto.PriceGood, error)
		GetStocks(ctx context.Context, apiKey string, warehouseID int64, req dto.GetStocksRequest) (*dto.GetStocksResponse, error)
		GetStocksByChrtIDs(ctx context.Context, apiKey string, warehouseID int64, chrtIDs []int64) ([]dto.StockItem, error)
		GetNewOrders(ctx context.Context, apiKey string) (*dto.GetNewOrdersResponse, error)
		GetNewDBWOrders(ctx context.Context, apiKey string) (*dto.GetNewDBWOrdersResponse, error)
		GetNewDBSOrders(ctx context.Context, apiKey string) (*dto.GetNewDBSOrdersResponse, error)
		GetSales(ctx context.Context, apiKey string, dateFrom time.Time, flag int, reportLoc *time.Location) (dto.GetSalesResponse, error)
		GetOrders(ctx context.Context, apiKey string, dateFrom time.Time, flag int, reportLoc *time.Location) (dto.GetOrdersResponse, error)
		GetWBWarehouses(ctx context.Context, apiKey string) ([]dto.WBWarehouse, error)
		GetSellerWarehouses(ctx context.Context, apiKey string) ([]dto.SellerWarehouse, error)
		GetFeedbacks(ctx context.Context, apiKey string, isAnswered bool, take, skip int, order string, dateFrom, dateTo *int64) (*dto.GetFeedbacksResponse, error)
		GetQuestions(ctx context.Context, apiKey string, isAnswered bool, take, skip int, order string, dateFrom, dateTo *int64) (*dto.GetQuestionsResponse, error)
		CheckToken(ctx context.Context, token string) error
	}

	TelegramAPI interface {
		SendMessage(ctx context.Context, chatID int64, text string) error
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
	fbsOrdersRepo FBSOrderRepository,
	dbwOrdersRepo DBWOrderRepository,
	dbsOrdersRepo DBSOrderRepository,
	saleEventsRepo SaleEventRepository,
	orderRefusalsRepo OrderRefusalRepository,
	wbOrdersRepo WBOrdersRepository,
	orderSyncCursorRepo OrderSyncCursorRepository,
	wbFeedbacksRepo WBFeedbackRepository,
	wbQuestionsRepo WBQuestionRepository,
	endpointsRepo EndpointsRepo,
	userRepo UserRepo,
	tg TelegramAPI,
	wb WBClient,
	defaultReportLoc *time.Location,
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
		fbsOrdersRepo:          fbsOrdersRepo,
		dbwOrdersRepo:          dbwOrdersRepo,
		dbsOrdersRepo:          dbsOrdersRepo,
		saleEventsRepo:         saleEventsRepo,
		orderRefusalsRepo:      orderRefusalsRepo,
		wbOrdersRepo:           wbOrdersRepo,
		orderSyncCursorRepo:    orderSyncCursorRepo,
		wbFeedbacksRepo:        wbFeedbacksRepo,
		wbQuestionsRepo:        wbQuestionsRepo,
		endpointsRepo:          endpointsRepo,
		userRepo:               userRepo,
		tg:                     tg,
		wbClient:               wb,
		defaultReportLoc:       defaultReportLoc,
		log:                    log,
	}
}

// shopReportLocation возвращает часовой пояс для отчётов WB по магазину: shop.WBReportTimezone или defaultReportLoc.
func (s *Service) shopReportLocation(shop *domain.Shop) *time.Location {
	if shop != nil && shop.WBReportTimezone != "" {
		loc, err := time.LoadLocation(shop.WBReportTimezone)
		if err == nil {
			return loc
		}
		s.log.Warn("shop wb_report_timezone load failed, using default", "shop_id", shop.ID, "tz", shop.WBReportTimezone, "err", err)
	}
	return s.defaultReportLoc
}

func wbStatusCode(err error) (int, bool) {
	var he *wbclient.HTTPError
	if errors.As(err, &he) {
		return he.StatusCode, true
	}
	return 0, false
}
