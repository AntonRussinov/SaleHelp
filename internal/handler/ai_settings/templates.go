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
// @Summary Сохранить шаблон AI по категории отзыва
// @Description Создаёт или обновляет шаблон для категории отзыва: дополнительные инструкции для AI и переопределение тональности.
// @Description
// @Description **Категории (`category` в пути):**
// @Description - `positive` — положительные отзывы (4–5 звёзд)
// @Description - `neutral` — нейтральные отзывы (3 звезды)
// @Description - `negative` — негативные отзывы (1–2 звезды)
// @Description
// @Description **Поля запроса:**
// @Description - `customPrompt` — дополнительные инструкции для AI (макс. 500 символов). Например: "Упомяни программу лояльности"
// @Description - `toneVoiceOverride` — переопределить тональность для этой категории: `standard` | `friendly` | `official` | `empathic` | `minimalist`. Пустая строка = использовать глобальную настройку
// @Description - `enabled` — включить/выключить шаблон для этой категории (по умолчанию `true`)
// @Tags ai-settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param category path string true "Категория отзыва" Enums(positive, neutral, negative)
// @Param request body api.UpsertAITemplateRequest true "Данные шаблона"
// @Success 200 {object} api.AITemplateResponse
// @Failure 400 {object} api.HTTPError "Невалидная категория или тело запроса"
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/ai-settings/templates/{category} [put]
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
