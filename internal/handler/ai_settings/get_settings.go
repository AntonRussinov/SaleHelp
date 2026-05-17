package ai_settings

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAISettings возвращает настройки AI магазина с правилами и шаблонами.
// Если настроек нет — возвращает дефолтные значения (не 404), чтобы фронт мог сразу показать форму.
func (h *Handler) GetAISettings(c echo.Context) error {
	shopID, shop, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	settings, err := h.repo.GetByShopID(c.Request().Context(), shopID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if settings == nil {
		empty := defaultSettings(shopID, shop.UserID)
		return c.JSON(http.StatusOK, settingsToResponse(empty, nil, nil))
	}

	rules, err := h.repo.ListRules(c.Request().Context(), settings.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	templates, err := h.repo.ListTemplates(c.Request().Context(), settings.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, settingsToResponse(settings, rules, templates))
}
