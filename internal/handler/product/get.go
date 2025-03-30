package product

import (
	"net/http"
	"salehelp/internal/handler"

	"github.com/labstack/echo/v4"
)

// GetUserProductByUserID
// @id           GetUserProductByUserID
// @Summary      Get user product data by userid
// @Description  Get user product data by userid
// @Security     sandstorm-backend
// @Tags         config, access
// @Produce      json
// @Param        id               path      string  true  "User ID"
// @Success      200                  {object}  api.GetUserProductResponse
// @Failure      400,401,403,404,500  {object}  api.HTTPError
// @Router       /v1/user-products/{id}  [get]
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
