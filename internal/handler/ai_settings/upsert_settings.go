package ai_settings

import (
	"net/http"

	"salehelp/api"
	"salehelp/domain"

	"github.com/labstack/echo/v4"
)

// UpsertAISettings создаёт или обновляет базовые настройки AI магазина.
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
