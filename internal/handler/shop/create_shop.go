package shop

import (
	"net/http"
	"strings"
	"time"

	"salehelp/api"
	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"
	"salehelp/pkg/wbtoken"

	"github.com/labstack/echo/v4"
)

// CreateShop создаёт магазин (декодирование и валидация WB-токена, дедупликация по sid).
//
// @Summary     Создать магазин
// @Description Подключение магазина WB: валидация JWT-токена, создание shop. Требует JWT.
// @Tags        shops
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body api.CreateShopRequest true "Название, токен WB, маркетплейс"
// @Success     201  {object} api.ShopResponse
// @Failure     400  {object} api.HTTPError
// @Failure     401  {object} api.HTTPError
// @Failure     409  {object} api.HTTPError
// @Router      /v1/shops [post]
func (h *Handler) CreateShop(c echo.Context) error {
	userID := auth_google_handler.UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req api.CreateShopRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token := strings.TrimSpace(req.APIToken)

	claims, err := wbtoken.Decode(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := wbtoken.Validate(claims); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Проверка лимита магазинов по тарифу.
	if err := h.checkShopLimit(c, userID); err != nil {
		return err
	}

	// Дедупликация по WB seller ID — один продавец может быть зарегистрирован только один раз.
	existing, err := h.shopsRepo.GetByExternalID(c.Request().Context(), claims.SID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ошибка проверки магазина"})
	}
	if existing != nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Этот магазин уже подключён к системе. Если это твой магазин — войди в тот аккаунт, к которому он привязан.",
		})
	}

	marketplace := strings.TrimSpace(strings.ToLower(req.Marketplace))
	if marketplace == "" {
		marketplace = "wildberries"
	}
	name := strings.TrimSpace(req.Name)
	if len(name) < 2 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name too short"})
	}

	tz := strings.TrimSpace(req.WBReportTimezone)
	if tz != "" {
		if _, err := time.LoadLocation(tz); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid timezone: " + tz})
		}
	}

	shop := &domain.Shop{
		UserID:           userID,
		Marketplace:      marketplace,
		ExternalID:       claims.SID,
		TokenType:        claims.TokenType(),
		Name:             name,
		APIToken:         token,
		IsActive:         true,
		WBReportTimezone: tz,
	}
	if err := h.shopCreator.CreateShop(c.Request().Context(), shop); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, shopToResponse(shop))
}
