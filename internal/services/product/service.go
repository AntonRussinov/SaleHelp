package product

import (
	"context"
	"fmt"
	"log/slog"
	"salehelp/domain"
	"salehelp/internal/adapter/kaspi/portal/models"
	"salehelp/pkg/logger/sl"
)

type (
	Repository interface {
		FindAllByMerchID(ctx context.Context, merchID string) ([]*domain.Product, error)
		// FindByID(ctx context.Context, id uuid.UUID) (*domain.Products, error)
		// Save(ctx context.Context, product *domain.Products) error
	}

	PortalAPI interface {
		GetAllProducts(ctx context.Context, merchantID string) (*models.Products, error)
	}

	Service struct {
		log       *slog.Logger
		repo      Repository
		portalAPI PortalAPI
	}
)

func NewService(log *slog.Logger, repo Repository, portalAPI PortalAPI) *Service {
	return &Service{
		log:       log.With(sl.Module("ProductsService")),
		repo:      repo,
		portalAPI: portalAPI,
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
