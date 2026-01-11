package subscription

import (
	"context"
	"salehelp/internal/services/subscription"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	Service interface {
		GrantSubscription(ctx context.Context, in GrantInput) (*subscription.GrantResult, error)
		GetSubscriptionStatus(ctx context.Context, shopID int) (*subscription.StatusResult, error)
	}

	Handler struct {
		service Service
	}
)

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (*Handler) getIDFromParam(c echo.Context) (int, error) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil || id <= 0 {
		return 0, echo.ErrBadRequest
	}

	return id, nil
}
