package alert

import (
	"context"
	"encoding/json"
	"net/http"
	"salehelp/api"
	"salehelp/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	Service interface {
		CreateAlert(ctx context.Context, shopID int, alert *domain.Alert) (*domain.Alert, error)
		UpdateAlert(ctx context.Context, alertID int, shopID int, alert *domain.Alert) (*domain.Alert, error)
		DeleteAlert(ctx context.Context, alertID int, shopID int) error
		GetAlert(ctx context.Context, alertID int, shopID int) (*domain.Alert, error)
		ListAlerts(ctx context.Context, shopID int, productID *int, alertType *string, status *string, limit, offset int) ([]*domain.Alert, int, error)
	}

	Handler struct {
		service Service
	}
)

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateAlert создает новый алерт
// @Summary Создать алерт
// @Description Создает новый алерт для магазина
// @Tags alerts
// @Accept json
// @Produce json
// @Param shopId path int true "ID магазина"
// @Param request body api.CreateAlertRequest true "Данные алерта"
// @Success 201 {object} api.AlertResponse
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v1/shops/{shopId}/alerts [post]
func (h *Handler) CreateAlert(c echo.Context) error {
	shopID, err := h.getShopIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
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

// UpdateAlert обновляет существующий алерт
// @Summary Обновить алерт
// @Description Обновляет существующий алерт
// @Tags alerts
// @Accept json
// @Produce json
// @Param shopId path int true "ID магазина"
// @Param alertId path int true "ID алерта"
// @Param request body api.UpdateAlertRequest true "Данные для обновления"
// @Success 200 {object} api.AlertResponse
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v1/shops/{shopId}/alerts/{alertId} [put]
func (h *Handler) UpdateAlert(c echo.Context) error {
	shopID, err := h.getShopIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
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

// DeleteAlert удаляет алерт
// @Summary Удалить алерт
// @Description Удаляет алерт по ID
// @Tags alerts
// @Param shopId path int true "ID магазина"
// @Param alertId path int true "ID алерта"
// @Success 204 "No Content"
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v1/shops/{shopId}/alerts/{alertId} [delete]
func (h *Handler) DeleteAlert(c echo.Context) error {
	shopID, err := h.getShopIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
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
// @Router /api/v1/shops/{shopId}/alerts/{alertId} [get]
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

// ListAlerts получает список алертов
// @Summary Получить список алертов
// @Description Получает список алертов с фильтрацией и пагинацией
// @Tags alerts
// @Produce json
// @Param shopId path int true "ID магазина"
// @Param productId query int false "ID товара для фильтрации"
// @Param type query string false "Тип алерта" Enums(price_drop, stock_low, stock_zero)
// @Param status query string false "Статус алерта" Enums(active, inactive)
// @Param limit query int false "Лимит записей (по умолчанию 20, максимум 100)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {object} api.ListAlertsResponse
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v1/shops/{shopId}/alerts [get]
func (h *Handler) ListAlerts(c echo.Context) error {
	shopID, err := h.getShopIDFromParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
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

// Вспомогательные методы

func (h *Handler) getShopIDFromParam(c echo.Context) (int, error) {
	paramID := c.Param("shopId")
	return strconv.Atoi(paramID)
}

func (h *Handler) getAlertIDFromParam(c echo.Context) (int, error) {
	paramID := c.Param("alertId")
	return strconv.Atoi(paramID)
}

func (h *Handler) domainAlertToResponse(alert *domain.Alert) api.AlertResponse {
	var payload map[string]interface{}
	if len(alert.Payload) > 0 {
		json.Unmarshal(alert.Payload, &payload)
	}

	return api.AlertResponse{
		ID:                 alert.ID,
		ShopID:             alert.ShopID,
		ProductID:          alert.ProductID,
		ProductSizeID:      alert.ProductSizeID,
		Type:               api.AlertType(alert.Type),
		Severity:           alert.Severity,
		Status:             api.AlertStatus(alert.Status),
		Message:            alert.Message,
		Payload:            payload,
		Scope:              api.AlertScope(alert.Scope),
		TriggeredAt:        alert.TriggeredAt,
		ResolvedAt:         alert.ResolvedAt,
		AIRecommendationID: alert.AIRecommendationID,
	}
}