package shop

import (
	"context"
	"strconv"

	"salehelp/domain"

	"github.com/labstack/echo/v4"
)

type (
	Presenter interface {
		// shop(shops *domain.shops) *view.shopView
	}

	Service interface {
		SyncShop(ctx context.Context, shopID int) error
	}

	ShopsRepo interface {
		GetByID(ctx context.Context, id int) (*domain.Shop, error)
		ListByUser(ctx context.Context, userID, limit, offset int) ([]*domain.Shop, error)
		Create(ctx context.Context, s *domain.Shop) error
	}

	ProductRepo interface {
		ListByShop(ctx context.Context, shopID, limit, offset int) ([]*domain.Product, error)
		CountByShop(ctx context.Context, shopID int) (int, error)
	}

	WBTokenChecker interface {
		CheckWBToken(ctx context.Context, token string) error
	}

	Handler struct {
		service    Service
		shopsRepo  ShopsRepo
		productRepo ProductRepo
		wbCheck    WBTokenChecker
	}
)

func NewHandler(service Service, shopsRepo ShopsRepo, productRepo ProductRepo, wbCheck WBTokenChecker) *Handler {
	return &Handler{
		service:    service,
		shopsRepo:  shopsRepo,
		productRepo: productRepo,
		wbCheck:    wbCheck,
	}
}

func (*Handler) getIDFromParam(c echo.Context) (int, error) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		return 0, echo.ErrBadRequest
	}

	return id, nil
}
