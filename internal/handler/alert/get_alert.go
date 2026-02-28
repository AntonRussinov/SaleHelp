package alert

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAlert получает алерт по ID
// @Summary Получить алерт
// @Description Получает алерт по ID
// @Tags alerts
// @Produce json
// @Param shopId path int true "ID магазина"
// @Param alertId path int true "ID алерта"
// @Success 200 {object} api.AlertResponse
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/alerts/{alertId} [get]
func (h *Handler) GetAlert(c echo.Context) error {
	shopID, err := h.getShopIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
	}

	alertID, err := h.getAlertIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid alert ID")
	}

	alert, err := h.service.GetAlert(c.Request().Context(), alertID, shopID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "alert not found")
	}

	response := h.domainAlertToResponse(alert)
	return c.JSON(http.StatusOK, response)
}
