package alert

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// DeleteAlert удаляет алерт
// @Summary Удалить алерт
// @Description Удаляет алерт по ID. Доступ только у владельца магазина (JWT).
// @Tags alerts
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param alertId path int true "ID алерта"
// @Success 204 "No Content"
// @Failure 400 {object} api.HTTPError
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/alerts/{alertId} [delete]
func (h *Handler) DeleteAlert(c echo.Context) error {
	shopID, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	alertID, err := h.getAlertIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid alert ID")
	}

	err = h.service.DeleteAlert(c.Request().Context(), alertID, shopID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "alert not found")
	}

	return c.NoContent(http.StatusNoContent)
}
