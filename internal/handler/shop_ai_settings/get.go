package shop_ai_settings

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAISettings возвращает настройки AI для магазина.
// Если настройки не заданы — 404.
//
// @Summary Получить настройки AI
// @Description Возвращает настройки AI-ответов на отзывы для магазина (бренд, тон, доп. правила)
// @Tags shop-ai-settings
// @Produce json
// @Param shopId path int true "ID магазина"
// @Success 200 {object} api.AISettingsResponse
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/ai-settings [get]
func (h *Handler) GetAISettings(c echo.Context) error {
	shopID, err := h.getShopIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
	}

	settings, err := h.repo.GetByShopID(c.Request().Context(), shopID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if settings == nil {
		return echo.NewHTTPError(http.StatusNotFound, "AI settings not found for this shop")
	}

	return c.JSON(http.StatusOK, h.toResponse(settings))
}
