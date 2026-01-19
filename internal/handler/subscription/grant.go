package subscription

import (
	"net/http"
	"strings"

	"salehelp/internal/handler"
	"salehelp/internal/services/subscription"

	"github.com/labstack/echo/v4"
)

// GrantSubscriptionByShopID
// @id           GrantSubscriptionByShopID
// @Summary      Grant/extend subscription for shop (manual v0)
// @Description  Creates a new subscription row (history-preserving) and upserts addons for this subscription.
// @Security     sandstorm-backend
// @Tags         subscription
// @Accept       json
// @Produce      json
// @Param        id   path      string true "Shop ID"
// @Param        body body      subscription.GrantInput true "Grant subscription payload"
// @Success      200  {object}  subscription.GrantResult
// @Failure      400,401,403,404,500 {object} api.HTTPError
// @Router       /v1/shop/{id}/subscription/grant [post]
func (hand *Handler) GrantSubscriptionByShopID(c echo.Context) error {
	const op = "handler.Subscription.GrantSubscriptionByShopID"

	shopID, err := hand.getIDFromParam(c)
	if err != nil {
		return echo.ErrNotFound
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
