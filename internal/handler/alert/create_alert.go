package alert

import (
	"encoding/json"
	"net/http"
	"salehelp/api"
	"salehelp/domain"

	"github.com/labstack/echo/v4"
)

// CreateAlert создает новый алерт
// @Summary Создать алерт
// @Description Создает новый алерт для магазина. Доступ только у владельца магазина (JWT).
// @Tags alerts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param request body api.CreateAlertRequest true "Данные алерта"
// @Success 201 {object} api.AlertResponse
// @Failure 400 {object} api.HTTPError
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/alerts [post]
func (h *Handler) CreateAlert(c echo.Context) error {
	shopID, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	var req api.CreateAlertRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.ValidatePayload(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload: "+err.Error())
	}

	// Конвертируем API модель в domain модель
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to marshal payload")
	}

	domainAlert := &domain.Alert{
		ShopID:        shopID,
		ProductID:     req.ProductID,
		ProductSizeID: req.ProductSizeID,
		Type:          string(req.Type),
		Severity:      req.Severity,
		Status:        string(req.Status),
		Message:       req.Message,
		Scope:         string(req.Scope),
		Payload:       payloadBytes,
	}

	createdAlert, err := h.service.CreateAlert(c.Request().Context(), shopID, domainAlert)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := h.domainAlertToResponse(createdAlert)
	return c.JSON(http.StatusCreated, response)
}
