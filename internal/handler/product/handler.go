package product

import (
	"context"
	"salehelp/domain"

	"github.com/labstack/echo/v4"
)

type (
	Presenter interface {
		// Product(products *domain.Products) *view.ProductView
	}

	Service interface {
		GetAllProducts(ctx context.Context, merchantID string) ([]*domain.Product, error)
		//GetAllProducts(ctx context.Context, merchantID string) (*domain.Products, error)
	}

	Handler struct {
		service   Service
		presenter Presenter
	}
)

func NewHandler(service Service, presenter Presenter) *Handler {
	return &Handler{
		service:   service,
		presenter: presenter,
	}
}

func (*Handler) getIDFromParam(c echo.Context) (string, error) {
	paramID := c.Param("id")

	return paramID, nil
}
