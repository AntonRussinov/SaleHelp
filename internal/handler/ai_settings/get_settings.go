package ai_settings

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAISettings возвращает настройки AI магазина с правилами и шаблонами.
// @Summary Получить AI-настройки магазина
// @Description Возвращает полные настройки AI автоответов на отзывы: основные параметры, правила бренда и шаблоны по категориям.
// @Description Если настройки ещё не созданы — возвращает дефолтные значения (не 404), чтобы фронт мог сразу показать форму.
// @Description
// @Description **Поля ответа:**
// @Description - `isActive` — включён ли AI автоответ для этого магазина
// @Description - `toneOfVoice` — тональность: `standard` | `friendly` | `official` | `empathic` | `minimalist`
// @Description - `responseLength` — длина ответа: `short` | `medium` | `long`
// @Description - `addressing` — обращение: `vy` (вы) | `ty` (ты)
// @Description - `emojiMode` — эмодзи: `none` | `minimal` | `expressive`
// @Description - `rules` — список правил бренда и запрещённых слов (см. `/ai-settings/rules`)
// @Description - `templates` — шаблоны по категориям отзывов: `positive` | `neutral` | `negative`
// @Tags ai-settings
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Success 200 {object} api.AISettingsResponse
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/ai-settings [get]
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
