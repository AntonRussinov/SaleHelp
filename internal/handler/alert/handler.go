package alert

import (
	"context"
	"encoding/json"
	"net/http"
	"salehelp/api"
	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ShopGetter нужен для проверки доступа: владелец магазина = user_id из JWT.
type ShopGetter interface {
	GetByID(ctx context.Context, id int) (*domain.Shop, error)
}

type (
	Service interface {
		CreateAlert(ctx context.Context, shopID int, alert *domain.Alert) (*domain.Alert, error)
		UpdateAlert(ctx context.Context, alertID int, shopID int, alert *domain.Alert) (*domain.Alert, error)
		DeleteAlert(ctx context.Context, alertID int, shopID int) error
		GetAlert(ctx context.Context, alertID int, shopID int) (*domain.Alert, error)
		ListAlerts(ctx context.Context, shopID int, productID *int, alertType *string, status *string, limit, offset int) ([]*domain.Alert, int, error)
	}

	Handler struct {
		service    Service
		shopGetter ShopGetter
	}
)

func NewHandler(service Service, shopGetter ShopGetter) *Handler {
	return &Handler{
		service:    service,
		shopGetter: shopGetter,
	}
}

// ensureShopAccess проверяет JWT (userID) и что пользователь — владелец магазина shopId. Возвращает ошибку 401/403.
func (h *Handler) ensureShopAccess(c echo.Context) (int, error) {
	userID := auth_google_handler.UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	shopID, err := h.getShopIDFromParam(c)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
	}
	shop, err := h.shopGetter.GetByID(c.Request().Context(), shopID)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if shop == nil || shop.UserID != userID {
		return 0, echo.NewHTTPError(http.StatusForbidden, "access denied")
	}
	return shopID, nil
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
