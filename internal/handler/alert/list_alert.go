package alert

import (
	"net/http"
	"salehelp/api"

	"github.com/labstack/echo/v4"
)

// ListAlerts получает список алертов
// @Summary Получить список алертов
// @Description Получает список алертов с фильтрацией и пагинацией. Доступ только у владельца магазина (JWT).
// @Tags alerts
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param productId query int false "ID товара для фильтрации"
// @Param type query string false "Тип алерта" Enums(price_drop, stock_low, stock_zero)
// @Param status query string false "Статус алерта" Enums(active, inactive)
// @Param limit query int false "Лимит записей (по умолчанию 20, максимум 100)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {object} api.ListAlertsResponse
// @Failure 400 {object} api.HTTPError
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/alerts [get]
func (h *Handler) ListAlerts(c echo.Context) error {
	shopID, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	var req api.ListAlertsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	// Значения по умолчанию
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var alertType *string
	if req.Type != nil {
		t := string(*req.Type)
		alertType = &t
	}

	var status *string
	if req.Status != nil {
		s := string(*req.Status)
		status = &s
	}

	alerts, total, err := h.service.ListAlerts(
		c.Request().Context(),
		shopID,
		req.ProductID,
		alertType,
		status,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := api.ListAlertsResponse{
		Alerts: make([]api.AlertResponse, len(alerts)),
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	for i, alert := range alerts {
		response.Alerts[i] = h.domainAlertToResponse(alert)
	}

	return c.JSON(http.StatusOK, response)
}
