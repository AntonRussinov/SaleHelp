package shop

import (
	"salehelp/internal/handler"

	"github.com/labstack/echo/v4"
)

// SyncShopByID запускает синхронизацию данных магазина с WB.
//
// @Summary     Синхронизация магазина
// @Description Запуск синка по ID магазина (товары, заказы, остатки и т.д.)
// @Tags        shops
// @Param       id path int true "ID магазина"
// @Success     204 "No Content"
// @Failure     400 {object} api.HTTPError
// @Failure     500 {object} api.HTTPError
// @Router      /v1/shop/{id}/sync [post]
func (hand *Handler) SyncShopByID(c echo.Context) error {
	const op = "handler.Usershop.SyncShopByID"

	id, err := hand.getIDFromParam(c)
	if err != nil {
		return echo.ErrNotFound
	}

	err = hand.service.SyncShop(c.Request().Context(), id)
	if err != nil {
		return handler.WrapError(op, err)
	}

	// return c.JSON(http.StatusOK, hand.presenter.shop(data))
	return c.NoContent(204)
}
