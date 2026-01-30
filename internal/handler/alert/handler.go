package alert

import (
	"context"
	"encoding/json"
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
