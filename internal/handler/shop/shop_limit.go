package shop

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// checkShopLimit возвращает ошибку если пользователь достиг лимита магазинов по тарифу.
// Если активной подписки нет — пропускает (первый магазин создаётся вместе с trial).
func (h *Handler) checkShopLimit(c echo.Context, userID int) error {
	ctx := c.Request().Context()

	sub, err := h.subsRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ошибка проверки тарифа"})
	}
	if sub == nil {
		// Нет подписки — первый магазин, trial создаётся в ShopCreator.
		return nil
	}

	count, err := h.shopsRepo.CountByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ошибка подсчёта магазинов"})
	}

	limit := sub.MaxShops + sub.ExtraShops
	if count >= limit {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Достигнут лимит магазинов для вашего тарифа",
			"code":  "shop_limit_reached",
		})
	}
	return nil
}
