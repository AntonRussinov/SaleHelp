package alert

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetAlert получает алерт по ID
// @Summary Получить алерт
// @Description Получает алерт по ID. Доступ только у владельца магазина (JWT).
// @Tags alerts
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param alertId path int true "ID алерта"
// @Success 200 {object} api.AlertResponse
// @Failure 400 {object} api.HTTPError
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/alerts/{alertId} [get]
func (h *Handler) GetAlert(c echo.Context) error {
	shopID, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	alertID, err := h.getAlertIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid alert ID")
	}

	alert, err := h.service.GetAlert(c.Request().Context(), alertID, shopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.log.Warn("GetAlert: not found",
				"request_id", requestID(c),
				"shop_id", shopID,
				"alert_id", alertID,
			)
			return echo.NewHTTPError(http.StatusNotFound, "alert not found")
		}
		h.log.Error("GetAlert: service error",
			"request_id", requestID(c),
			"shop_id", shopID,
			"alert_id", alertID,
			"error", err,
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := h.domainAlertToResponse(alert)
	return c.JSON(http.StatusOK, response)
}
