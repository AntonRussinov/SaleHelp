package product

import (
	"net/http"
	"salehelp/internal/handler"

	"github.com/labstack/echo/v4"
)

// GetProductByMerchantID возвращает список товаров (по merchant/shop ID из query).
//
// @Summary     Список товаров
// @Description Товары по идентификатору (query). Используйте GET /v1/shops/{shopId}/products для списка по магазину с JWT.
// @Tags        products
// @Produce     json
// @Param       id query string false "Идентификатор магазина/мерчанта"
// @Success     200 {array} domain.Product
// @Failure     400 {object} api.HTTPError
// @Failure     500 {object} api.HTTPError
// @Router      /v1/products [get]
func (hand *Handler) GetProductByMerchantID(c echo.Context) error {
	const op = "handler.UserProduct.GetUserProductByUserID"

	id, err := hand.getIDFromParam(c)
	if err != nil {
		return echo.ErrNotFound
	}

	data, err := hand.service.GetAllProducts(c.Request().Context(), id)
	if err != nil {
		return handler.WrapError(op, err)
	}

	// return c.JSON(http.StatusOK, hand.presenter.Product(data))
	return c.JSON(http.StatusOK, data)
}
