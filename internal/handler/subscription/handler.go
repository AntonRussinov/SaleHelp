package subscription

import (
	"context"
	"net/http"
	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"
	"salehelp/internal/services/subscription"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ShopGetter для проверки доступа: владелец магазина = user_id из JWT.
type ShopGetter interface {
	GetByID(ctx context.Context, id int) (*domain.Shop, error)
}

type (
	Service interface {
		GrantSubscription(ctx context.Context, in subscription.GrantInput) (*subscription.GrantResult, error)
		GetSubscriptionStatus(ctx context.Context, shopID int) (*subscription.StatusResult, error)
	}

	Handler struct {
		service    Service
		shopGetter ShopGetter
	}
)

func NewHandler(service Service, shopGetter ShopGetter) *Handler {
	return &Handler{
		service:    service,
		shopGetter: shopGetter,
	}
}

// ensureShopAccessByID проверяет JWT и что пользователь — владелец магазина (параметр "id"). Для роутов /shop/:id/...
func (hand *Handler) ensureShopAccessByID(c echo.Context) (int, error) {
	userID := auth_google_handler.UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return 0, c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	shopID, err := hand.getIDFromParam(c)
	if err != nil {
		return 0, err
	}
	shop, err := hand.shopGetter.GetByID(c.Request().Context(), shopID)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if shop == nil || shop.UserID != userID {
		return 0, c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}
	return shopID, nil
}

func (*Handler) getIDFromParam(c echo.Context) (int, error) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil || id <= 0 {
		return 0, echo.ErrBadRequest
	}

	return id, nil
}
