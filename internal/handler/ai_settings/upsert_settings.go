package ai_settings

import (
	"net/http"

	"salehelp/api"
	"salehelp/domain"

	"github.com/labstack/echo/v4"
)

// UpsertAISettings создаёт или обновляет базовые настройки AI магазина.
// @Summary Сохранить AI-настройки магазина
// @Description Создаёт или обновляет основные параметры AI автоответов для магазина.
// @Description Возвращает полный объект настроек включая существующие правила и шаблоны.
// @Description
// @Description **Допустимые значения полей:**
// @Description - `toneOfVoice`: `standard` | `friendly` | `official` | `empathic` | `minimalist`
// @Description - `responseLength`: `short` | `medium` | `long`
// @Description - `addressing`: `vy` (Вы) | `ty` (ты)
// @Description - `emojiMode`: `none` | `minimal` | `expressive`
// @Description - `brandName` — макс. 100 символов
// @Description - `brandDescription` — макс. 500 символов
// @Description - `representativeName` — имя представителя в ответах, макс. 50 символов
// @Description - `signature` — подпись в конце ответа, макс. 200 символов
// @Tags ai-settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param request body api.UpsertAISettingsRequest true "Настройки AI"
// @Success 200 {object} api.AISettingsResponse
// @Failure 400 {object} api.HTTPError "Невалидное тело запроса или недопустимое значение поля"
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/ai-settings [put]
func (h *Handler) UpsertAISettings(c echo.Context) error {
	shopID, shop, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	var req api.UpsertAISettingsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existing, err := h.repo.GetByShopID(c.Request().Context(), shopID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	d := &domain.AISettings{
		ShopID:             shopID,
		UserID:             shop.UserID,
		IsActive:           req.IsActive,
		ToneOfVoice:        req.ToneOfVoice,
		ResponseLength:     req.ResponseLength,
		Addressing:         req.Addressing,
		EmojiMode:          req.EmojiMode,
		BrandName:          req.BrandName,
		BrandDescription:   req.BrandDescription,
		RepresentativeName: req.RepresentativeName,
		Signature:          req.Signature,
	}
	if existing != nil {
		d.ID = existing.ID
		d.CreatedAt = existing.CreatedAt
	}

	if err := h.repo.Upsert(c.Request().Context(), d); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	rules, err := h.repo.ListRules(c.Request().Context(), d.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	templates, err := h.repo.ListTemplates(c.Request().Context(), d.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, settingsToResponse(d, rules, templates))
}
