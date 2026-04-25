package shop

import (
	"net/http"
	"strings"
	"time"

	"salehelp/api"
	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"

	"github.com/labstack/echo/v4"
)

// CreateShop создаёт магазин (проверка WB-токена, создание записи).
//
// @Summary     Создать магазин
// @Description Подключение магазина WB: проверка apiToken, создание shop. Требует JWT.
// @Tags        shops
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body api.CreateShopRequest true "Название, токен WB, маркетплейс"
// @Success     201  {object} api.ShopResponse
// @Failure     400  {object} api.HTTPError
// @Failure     401  {object} api.HTTPError
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
	if len(token) < 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "apiToken too short"})
	}
	if err := h.wbCheck.CheckWBToken(c.Request().Context(), token); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid wb token: " + err.Error()})
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

	tokenType := strings.TrimSpace(req.TokenType)
	if tokenType != domain.TokenTypeService {
		tokenType = domain.TokenTypeBasic
	}

	shop := &domain.Shop{
		UserID:           userID,
		Marketplace:      marketplace,
		Name:             name,
		APIToken:         token,
		TokenType:        tokenType,
		IsActive:         true,
		WBReportTimezone: tz,
	}
	if err := h.shopsRepo.Create(c.Request().Context(), shop); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, shopToResponse(shop))
}
