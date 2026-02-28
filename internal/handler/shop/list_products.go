package shop

import (
	"net/http"
	"strconv"

	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"

	"github.com/labstack/echo/v4"
)

// ProductResponse продукт в ответе GET /shops/:shopId/products.
type ProductResponse struct {
	ID         int    `json:"id"`
	ShopID     int    `json:"shopId"`
	ExternalID string `json:"externalId"`
	Title      string `json:"title"`
	Brand      string `json:"brand,omitempty"`
}

// ListShopProducts GET /shops/:shopId/products — список товаров магазина (JWT, проверка доступа).
func (h *Handler) ListShopProducts(c echo.Context) error {
	userID := auth_google_handler.UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	shopID, err := getShopIDFromParam(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid shopId"})
	}

	shop, err := h.shopsRepo.GetByID(c.Request().Context(), shopID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if shop == nil || shop.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	products, err := h.productRepo.ListByShop(c.Request().Context(), shopID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	total, _ := h.productRepo.CountByShop(c.Request().Context(), shopID)

	out := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		out = append(out, productToResponse(p))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"items":  out,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func getShopIDFromParam(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("shopId"))
	if err != nil || id <= 0 {
		return 0, echo.ErrBadRequest
	}
	return id, nil
}

func productToResponse(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:         p.ID,
		ShopID:     p.ShopID,
		ExternalID: p.ExternalID,
		Title:      p.Title,
		Brand:      p.Brand,
	}
}
