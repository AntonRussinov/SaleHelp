package ai_settings

import (
	"net/http"

	"salehelp/api"
	"salehelp/domain"

	"github.com/labstack/echo/v4"
)

var validCategories = map[string]bool{
	domain.AITemplateCategoryPositive: true,
	domain.AITemplateCategoryNeutral:  true,
	domain.AITemplateCategoryNegative: true,
}

// UpsertTemplate создаёт или обновляет шаблон для конкретной категории отзыва.
func (h *Handler) UpsertTemplate(c echo.Context) error {
	shopID, shop, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	category := c.Param("category")
	if !validCategories[category] {
		return echo.NewHTTPError(http.StatusBadRequest, "category must be positive, neutral or negative")
	}

	var req api.UpsertAITemplateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	settings, err := h.ensureSettings(c, shopID, shop.UserID)
	if err != nil {
		return err
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	tmpl := &domain.AITemplate{
		AISettingsID:      settings.ID,
		Category:          category,
		CustomPrompt:      req.CustomPrompt,
		ToneVoiceOverride: req.ToneVoiceOverride,
		Enabled:           enabled,
	}
	if err := h.repo.UpsertTemplate(c.Request().Context(), tmpl); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, api.AITemplateResponse{
		ID:                tmpl.ID,
		Category:          tmpl.Category,
		CustomPrompt:      tmpl.CustomPrompt,
		ToneVoiceOverride: tmpl.ToneVoiceOverride,
		Enabled:           tmpl.Enabled,
	})
}
