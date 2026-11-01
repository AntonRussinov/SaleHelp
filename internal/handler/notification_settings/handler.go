package notification_settings

import (
	"context"
	"net/http"
	"salehelp/api"
	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"
	"salehelp/internal/notificationevents"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ShopGetter interface {
	GetByID(ctx context.Context, id int) (*domain.Shop, error)
}

type Service interface {
	EffectiveMap(ctx context.Context, shopID int) (map[string]bool, error)
	Set(ctx context.Context, shopID int, eventKey string, enabled bool) error
}

type Handler struct {
	svc        Service
	shopGetter ShopGetter
}

func NewHandler(svc Service, shopGetter ShopGetter) *Handler {
	return &Handler{svc: svc, shopGetter: shopGetter}
}

// ensureShopAccess проверяет JWT (userID) и что пользователь — владелец магазина shopId. Возвращает ошибку 401/403.
func (h *Handler) ensureShopAccess(c echo.Context) (int, error) {
	userID := auth_google_handler.UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return 0, c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	shopID, err := strconv.Atoi(c.Param("shopId"))
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
	}
	shop, err := h.shopGetter.GetByID(c.Request().Context(), shopID)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if shop == nil || shop.UserID != userID {
		return 0, c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}
	return shopID, nil
}

// ListSettings
// @Summary Список настроек уведомлений магазина
// @Description Возвращает список всех доступных типов уведомлений и их текущее состояние (enabled/disabled) для магазина.
// @Tags notification-settings
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Success 200 {object} api.ListNotificationSettingsResponse
// @Failure 400 {object} api.HTTPError
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/notification-settings [get]
func (h *Handler) ListSettings(c echo.Context) error {
	shopID, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	m, err := h.svc.EffectiveMap(c.Request().Context(), shopID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	items := make([]api.NotificationSetting, 0, len(notificationevents.Catalog))
	for _, def := range notificationevents.All() {
		enabled := def.DefaultEnabled
		if v, ok := m[def.Key]; ok {
			enabled = v
		}
		items = append(items, api.NotificationSetting{
			Key:            def.Key,
			Enabled:        enabled,
			Title:          def.Title,
			Description:    def.Description,
			Category:       def.Category,
			DefaultEnabled: def.DefaultEnabled,
		})
	}

	return c.JSON(http.StatusOK, api.ListNotificationSettingsResponse{Items: items})
}

// UpdateSettings
// @Summary Обновить настройки уведомлений магазина
// @Description Обновляет enabled/disabled для указанных типов уведомлений (владелец магазина).
// @Tags notification-settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param shopId path int true "ID магазина"
// @Param request body api.UpdateNotificationSettingsRequest true "Список настроек"
// @Success 200 {object} api.ListNotificationSettingsResponse
// @Failure 400 {object} api.HTTPError
// @Failure 401 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /v1/shops/{shopId}/notification-settings [put]
func (h *Handler) UpdateSettings(c echo.Context) error {
	shopID, err := h.ensureShopAccess(c)
	if err != nil {
		return err
	}

	var req api.UpdateNotificationSettingsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, it := range req.Items {
		if _, ok := notificationevents.Find(it.Key); !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "unknown notification key: "+it.Key)
		}
		if err := h.svc.Set(c.Request().Context(), shopID, it.Key, it.Enabled); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// Return updated list.
	return h.ListSettings(c)
}
