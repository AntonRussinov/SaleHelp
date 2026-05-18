package ai_settings

import (
	"net/http"
	"strconv"

	"salehelp/api"
	"salehelp/domain"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreateRule добавляет правило бренда или запрещённое слово.
// @Summary Добавить правило AI
// @Description Добавляет правило бренда или запрещённое слово для AI автоответов.
// @Description Лимит правил зависит от тарифа пользователя.
// @Description
// @Description **Типы правил (`ruleType`):**
// @Description - `rule` — инструкция для AI (например: "Всегда упоминай гарантию 1 год")
// @Description - `forbidden_word` — слово/фраза, которую AI не должен использовать
// @Description
// @Description `orderIndex` — порядок применения правил (0 = первое).
// @Tags ai-settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param request body api.CreateAIRuleRequest true "Данные правила"
// @Success 201 {object} api.AIRuleResponse
// @Failure 400 {object} api.HTTPError "Невалидное тело или недопустимый тип правила"
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError "Нет доступа к магазину или достигнут лимит правил тарифа"
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/ai-settings/rules [post]
func (h *Handler) CreateRule(c echo.Context) error {
	shopID, shop, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	var req api.CreateAIRuleRequest
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

	// Проверяем лимит правил по тарифу.
	if err := h.checkRuleLimit(c, shop.UserID, settings.ID); err != nil {
		return err
	}

	rule := &domain.AIRule{
		AISettingsID: settings.ID,
		RuleType:     req.RuleType,
		RuleText:     req.RuleText,
		OrderIndex:   req.OrderIndex,
	}
	if err := h.repo.CreateRule(c.Request().Context(), rule); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, api.AIRuleResponse{
		ID:         rule.ID,
		RuleType:   rule.RuleType,
		RuleText:   rule.RuleText,
		OrderIndex: rule.OrderIndex,
	})
}

// UpdateRule обновляет правило.
// @Summary Обновить правило AI
// @Description Обновляет текст, тип и порядок существующего правила AI.
// @Description Правило должно принадлежать настройкам этого магазина.
// @Tags ai-settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param ruleId path int true "ID правила"
// @Param request body api.UpdateAIRuleRequest true "Обновлённые данные правила"
// @Success 200 {object} api.AIRuleResponse
// @Failure 400 {object} api.HTTPError "Невалидное тело или недопустимый ruleId"
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError "AI-настройки магазина не найдены"
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/ai-settings/rules/{ruleId} [put]
func (h *Handler) UpdateRule(c echo.Context) error {
	shopID, _, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	ruleID, err := strconv.Atoi(c.Param("ruleId"))
	if err != nil || ruleID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid rule ID")
	}

	var req api.UpdateAIRuleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	settings, err := h.repo.GetByShopID(c.Request().Context(), shopID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if settings == nil {
		return echo.NewHTTPError(http.StatusNotFound, "AI settings not found")
	}

	rule := &domain.AIRule{
		ID:           ruleID,
		AISettingsID: settings.ID,
		RuleType:     req.RuleType,
		RuleText:     req.RuleText,
		OrderIndex:   req.OrderIndex,
	}
	if err := h.repo.UpdateRule(c.Request().Context(), rule); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, api.AIRuleResponse{
		ID:         rule.ID,
		RuleType:   rule.RuleType,
		RuleText:   rule.RuleText,
		OrderIndex: rule.OrderIndex,
	})
}

// DeleteRule удаляет правило.
// @Summary Удалить правило AI
// @Description Удаляет правило бренда или запрещённое слово по его ID.
// @Tags ai-settings
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param ruleId path int true "ID правила"
// @Success 204 "Правило удалено"
// @Failure 400 {object} api.HTTPError "Невалидный ruleId"
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError "Правило или AI-настройки не найдены"
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/ai-settings/rules/{ruleId} [delete]
func (h *Handler) DeleteRule(c echo.Context) error {
	shopID, _, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	ruleID, err := strconv.Atoi(c.Param("ruleId"))
	if err != nil || ruleID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid rule ID")
	}

	settings, err := h.repo.GetByShopID(c.Request().Context(), shopID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if settings == nil {
		return echo.NewHTTPError(http.StatusNotFound, "AI settings not found")
	}

	if err := h.repo.DeleteRule(c.Request().Context(), settings.ID, ruleID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "rule not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// ensureSettings возвращает существующие настройки или создаёт дефолтные.
func (h *Handler) ensureSettings(c echo.Context, shopID, userID int) (*domain.AISettings, error) {
	settings, err := h.repo.GetByShopID(c.Request().Context(), shopID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if settings != nil {
		return settings, nil
	}
	d := defaultSettings(shopID, userID)
	if err := h.repo.Upsert(c.Request().Context(), d); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return d, nil
}

// checkRuleLimit проверяет что пользователь не превысил лимит правил по тарифу.
func (h *Handler) checkRuleLimit(c echo.Context, userID, settingsID int) error {
	sub, err := h.subsRepo.GetActiveByUserID(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tier := domain.TierTrial
	if sub != nil {
		tier = sub.Tier
	}
	maxRules := maxRulesForTier(tier)
	if maxRules == 0 {
		return echo.NewHTTPError(http.StatusForbidden, map[string]string{
			"error": "Правила бренда недоступны на вашем тарифе",
			"code":  "rules_not_available",
		})
	}

	existing, err := h.repo.ListRules(c.Request().Context(), settingsID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if len(existing) >= maxRules {
		return echo.NewHTTPError(http.StatusForbidden, map[string]string{
			"error": "Достигнут лимит правил для вашего тарифа",
			"code":  "rules_limit_reached",
		})
	}
	return nil
}
