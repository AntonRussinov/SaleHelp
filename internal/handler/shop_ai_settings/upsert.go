package shop_ai_settings

import (
	"net/http"
	"salehelp/api"
	"salehelp/domain"

	"github.com/labstack/echo/v4"
)

// UpsertAISettings создаёт или обновляет настройки AI для магазина.
// toneOfVoice проверяется на допустимые значения: standard | friendly | official | empathic | minimalist.
//
// @Summary Создать/обновить настройки AI
// @Description Задаёт настройки для формирования промпта ИИ: бренд, тон ответа, доп. правила, вкл/выкл генерации
// @Tags shop-ai-settings
// @Accept json
// @Produce json
// @Param shopId path int true "ID магазина"
// @Param request body api.UpsertAISettingsRequest true "Настройки AI"
// @Success 200 {object} api.AISettingsResponse
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v1/shops/{shopId}/ai-settings [put]
func (h *Handler) UpsertAISettings(c echo.Context) error {
	shopID, err := h.getShopIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
	}

	var req api.UpsertAISettingsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existing, _ := h.repo.GetByShopID(c.Request().Context(), shopID)
	d := &domain.ShopAISettings{
		ShopID:      shopID,
		ToneOfVoice: req.ToneOfVoice,
		BrandName:   req.BrandName,
		CustomRules: req.CustomRules,
		IsActive:    req.IsActive,
	}
	if existing != nil {
		d.CreatedAt = existing.CreatedAt
	}

	if err := h.repo.Save(c.Request().Context(), d); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Возвращаем сохранённые данные с актуальными датами из БД
	saved, _ := h.repo.GetByShopID(c.Request().Context(), shopID)
	if saved != nil {
		return c.JSON(http.StatusOK, h.toResponse(saved))
	}
	return c.JSON(http.StatusOK, h.toResponse(d))
}
