package shop

import (
	"net/http"
	"strconv"

	"salehelp/api"
	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"

	"github.com/labstack/echo/v4"
)

// ListShops возвращает список магазинов текущего пользователя.
//
// @Summary     Список магазинов
// @Description Магазины пользователя (JWT). Query: limit, offset.
// @Tags        shops
// @Produce     json
// @Security    BearerAuth
// @Param       limit  query int false "Лимит" default(20)
// @Param       offset query int false "Смещение" default(0)
// @Success     200 {array} api.ShopResponse
// @Failure     401 {object} api.HTTPError
// @Router      /v1/shops [get]
func (h *Handler) ListShops(c echo.Context) error {
	userID := auth_google_handler.UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	ctx := c.Request().Context()
	shops, err := h.shopsRepo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	out := make([]api.ShopResponse, 0, len(shops))
	for _, s := range shops {
		out = append(out, shopToResponse(s))
	}
	return c.JSON(http.StatusOK, out)
}

func shopToResponse(s *domain.Shop) api.ShopResponse {
	return api.ShopResponse{
		ID:                 s.ID,
		UserID:             s.UserID,
		Marketplace:        s.Marketplace,
		Name:               s.Name,
		IsActive:           s.IsActive,
		TokenType:          s.TokenType,
		WBReportTimezone:   s.WBReportTimezone,
		LastSyncAt:         s.LastSyncAt,
		LastSyncStatusCode: s.LastSyncStatusCode,
		CreatedAt:          s.CreatedAt,
	}
}
