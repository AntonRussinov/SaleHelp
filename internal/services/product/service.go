package product

import (
	"context"
	"fmt"
	"log/slog"
	"salehelp/domain"
	"salehelp/internal/adapter/kaspi/portal/portal_models"
	"salehelp/pkg/logger/sl"
)

type (
	Repository interface {
		FindAllByMerchID(ctx context.Context, merchID string) ([]*domain.Product, error)
		// FindByID(ctx context.Context, id uuid.UUID) (*domain.Products, error)
		// Save(ctx context.Context, product *domain.Products) error
	}

	StockRepo interface {
		Save(ctx context.Context, stock *domain.Stock) error
	}

	ProductRepo interface {
		Save(ctx context.Context, product *domain.Product) error
	}

	OfferRepo interface {
		Save(ctx context.Context, offer *domain.Offer) error
	}

	MerchantRepo interface {
		Save(ctx context.Context, merchant *domain.Merchant) error
	}

	ImageRepo interface {
		Save(ctx context.Context, iamge *domain.Image) error
	}

	CityPriceRepo interface {
		Save(ctx context.Context, cityPrice *domain.CityPrice) error
	}
	PortalAPI interface {
		GetAllProducts(ctx context.Context, merchantID string) (*portal_models.Products, error)
	}

	Service struct {
		log           *slog.Logger
		StockRepo     StockRepo
		ProductRepo   ProductRepo
		CityPriceRepo CityPriceRepo
		OfferRepo     OfferRepo
		MerchantRepo  MerchantRepo
		ImageRepo     ImageRepo
		portalAPI     PortalAPI
	}
)

func NewService(log *slog.Logger,
	StockRepo StockRepo, ProductRepo ProductRepo,
	CityPriceRepo CityPriceRepo, OfferRepo OfferRepo,
	MerchantRepo MerchantRepo, ImageRepo ImageRepo,
	portalAPI PortalAPI,
) *Service {
	return &Service{
		log:           log.With(sl.Module("ProductsService")),
		StockRepo:     StockRepo,
		ProductRepo:   ProductRepo,
		CityPriceRepo: CityPriceRepo,
		OfferRepo:     OfferRepo,
		MerchantRepo:  MerchantRepo,
		ImageRepo:     ImageRepo,
		portalAPI:     portalAPI,
	}
}

func (service *Service) handleError(err error, fn, message string) error {
	return domain.NewUnexpectedErr(fmt.Sprintf("%s: %s", fn, message), err)
}

func notFoundError(id string, err error) error {
	return domain.NewNotFoundWithCauseErr("product", "id", id, err)
}

func productRequestNotFoundError(id string, err error) error {
	return domain.NewNotFoundWithCauseErr("ProductsRequest", "user_id", id, err)
}
