package shop

import (
	"salehelp/internal/handler"

	"github.com/labstack/echo/v4"
)

// GetUsershopByUserID
// @id           GetUsershopByUserID
// @Summary      Get user shop data by userid
// @Description  Get user shop data by userid
// @Security     sandstorm-backend
// @Tags         config, access
// @Produce      json
// @Param        id               path      string  true  "User ID"
// @Success      204                  {object}
// @Failure      400,401,403,404,500  {object}  api.HTTPError
// @Router       /v1/shop/{id}/sync  [post]
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
