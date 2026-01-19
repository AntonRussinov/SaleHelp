package subscription

import (
	"net/http"

	"salehelp/internal/handler"

	"github.com/labstack/echo/v4"
)

// GetSubscriptionStatusByShopID
// @id           GetSubscriptionStatusByShopID
// @Summary      Get subscription status & entitlements for shop
// @Description  Returns active subscription (if any) and telegram chats limits/usage.
// @Security     sandstorm-backend
// @Tags         subscription
// @Produce      json
// @Param        id   path      string true "Shop ID"
// @Success      200  {object}  subscription.StatusResult
// @Failure      400,401,403,404,500 {object} api.HTTPError
// @Router       /v1/shop/{id}/subscription/status [get]
func (hand *Handler) GetSubscriptionStatusByShopID(c echo.Context) error {
	const op = "handler.Subscription.GetSubscriptionStatusByShopID"

	shopID, err := hand.getIDFromParam(c)
	if err != nil {
		return echo.ErrNotFound
	}

	res, err := hand.service.GetSubscriptionStatus(c.Request().Context(), shopID)
	if err != nil {
		return handler.WrapError(op, err)
	}

	return c.JSON(http.StatusOK, res)
}
