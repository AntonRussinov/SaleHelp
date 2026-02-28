package subscription

import (
	"net/http"
	"strings"

	"salehelp/internal/handler"
	"salehelp/internal/services/subscription"

	"github.com/labstack/echo/v4"
)

// GrantSubscriptionByShopID выдать или продлить подписку по магазину. Доступ только у пользователя с ролью admin (JWT).
//
// @Summary     Выдать/продлить подписку
// @Description Создаёт запись подписки (period_months: 1|3|6), опционально аддоны на доп. Telegram-чаты. Вызывать может только админ.
// @Tags        subscription
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path int true "ID магазина"
// @Param       body body subscription.GrantInput true "Период, валюта, суммы, доп. чаты"
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

	// path param is source of truth
	req.ShopID = shopID
	req.CurrencyCode = strings.ToUpper(strings.TrimSpace(req.CurrencyCode))
	req.Source = strings.TrimSpace(req.Source)

	// минимальная валидация v0
	if req.PeriodMonths != 1 && req.PeriodMonths != 3 && req.PeriodMonths != 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "period_months must be 1, 3 or 6")
	}
	if len(req.CurrencyCode) != 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "currency_code must be 3-letter ISO code (e.g. RUB, KZT)")
	}
	if req.BaseAmountMinor < 0 || req.TelegramAddonAmountMinor < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "amounts must be >= 0")
	}
	if req.ExtraTelegramChats < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "extra_telegram_chats must be >= 0")
	}
	if req.Source == "" {
		req.Source = "manual"
	}

	res, err := hand.service.GrantSubscription(c.Request().Context(), req)
	if err != nil {
		return handler.WrapError(op, err)
	}

	return c.JSON(http.StatusOK, res)
}
