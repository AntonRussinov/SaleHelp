package subscription

import (
	"net/http"

	"salehelp/internal/handler"
	"salehelp/internal/services/subscription"

	"github.com/labstack/echo/v4"
)

// GrantSubscriptionByShopID выдать подписку по магазину. Доступ только у пользователя с ролью admin (JWT).
//
// @Summary     Выдать/продлить подписку
// @Description Создаёт подписку для пользователя, владеющего магазином. Вызывать может только админ.
// @Tags        subscription
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path int true "ID магазина"
// @Param       body body subscription.GrantInput true "Тариф, период, сумма"
// @Success     200  {object} subscription.GrantResult
// @Failure     400  {object} api.HTTPError
// @Failure     401  {object} api.HTTPError
// @Failure     403  {object} api.HTTPError
// @Failure     500  {object} api.HTTPError
// @Router      /v1/shop/{id}/subscription/grant [post]
func (hand *Handler) GrantSubscriptionByShopID(c echo.Context) error {
	const op = "handler.Subscription.GrantSubscriptionByShopID"

	shopID, err := hand.getIDFromParam(c)
	if err != nil {
		return err
	}

	var req subscription.GrantInput
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	req.ShopID = shopID

	if req.Tier == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "tier is required")
	}
	if req.PeriodMonths < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "period_months must be >= 1")
	}

	res, err := hand.service.GrantSubscription(c.Request().Context(), req)
	if err != nil {
		return handler.WrapError(op, err)
	}

	return c.JSON(http.StatusOK, res)
}
