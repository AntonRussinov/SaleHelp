package subscription

import (
	"net/http"

	"salehelp/internal/handler"

	"github.com/labstack/echo/v4"
)

// GetSubscriptionStatusByShopID возвращает статус подписки по магазину.
//
// @Summary     Статус подписки
// @Description Активная подписка (если есть), срок, лимиты и использование Telegram-чатов
// @Tags        subscription
// @Produce     json
// @Param       id path int true "ID магазина"
// @Success     200 {object} subscription.StatusResult
// @Failure     400 {object} api.HTTPError
// @Failure     500 {object} api.HTTPError
// @Router      /v1/shop/{id}/subscription/status [get]
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
