package alert

import (
	"encoding/json"
	"net/http"
	"salehelp/api"

	"github.com/labstack/echo/v4"
)

// UpdateAlert обновляет существующий алерт
// @Summary Обновить алерт
// @Description Обновляет существующий алерт. Доступ только у владельца магазина (JWT).
// @Tags alerts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param alertId path int true "ID алерта"
// @Param request body api.UpdateAlertRequest true "Данные для обновления"
// @Success 200 {object} api.AlertResponse
// @Failure 400 {object} api.HTTPError
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/alerts/{alertId} [put]
func (h *Handler) UpdateAlert(c echo.Context) error {
	shopID, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	alertID, err := h.getAlertIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid alert ID")
	}

	var req api.UpdateAlertRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Получаем существующий алерт
	existingAlert, err := h.service.GetAlert(c.Request().Context(), alertID, shopID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "alert not found")
	}

	// Обновляем только переданные поля
	updatedAlert := *existingAlert

	if req.ProductID != nil {
		updatedAlert.ProductID = req.ProductID
	}
	if req.ProductSizeID != nil {
		updatedAlert.ProductSizeID = req.ProductSizeID
	}
	if req.Severity != nil {
		updatedAlert.Severity = *req.Severity
	}
	if req.Status != nil {
		updatedAlert.Status = string(*req.Status)
	}
	if req.Message != nil {
		updatedAlert.Message = *req.Message
	}
	if req.Scope != nil {
		updatedAlert.Scope = string(*req.Scope)
	}
	if req.Payload != nil {
		payloadBytes, err := json.Marshal(req.Payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to marshal payload")
		}
		updatedAlert.Payload = payloadBytes
	}

	result, err := h.service.UpdateAlert(c.Request().Context(), alertID, shopID, &updatedAlert)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := h.domainAlertToResponse(result)
	return c.JSON(http.StatusOK, response)
}
