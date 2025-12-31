package shop

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	Presenter interface {
		// shop(shops *domain.shops) *view.shopView
	}

	Service interface {
		SyncShop(ctx context.Context, shopID int) error
	}

	Handler struct {
		service Service
		// presenter Presenter
	}
)

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
		// presenter: presenter,
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
