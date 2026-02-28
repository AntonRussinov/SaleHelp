package shop_ai_settings

import (
	"context"
	"net/http"
	"salehelp/api"
	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"
	"strconv"

	"github.com/labstack/echo/v4"
)

// AISettingsRepo интерфейс репозитория настроек AI.
type AISettingsRepo interface {
	GetByShopID(ctx context.Context, shopID int) (*domain.ShopAISettings, error)
	Save(ctx context.Context, d *domain.ShopAISettings) error
}

// ShopGetter для проверки доступа: владелец магазина = user_id из JWT.
type ShopGetter interface {
	GetByID(ctx context.Context, id int) (*domain.Shop, error)
}

// Handler HTTP-хендлеры для настроек AI магазина (бренд, тон, правила для промпта).
type Handler struct {
	repo       AISettingsRepo
	shopGetter ShopGetter
}

// NewHandler создаёт хендлер.
func NewHandler(repo AISettingsRepo, shopGetter ShopGetter) *Handler {
	return &Handler{repo: repo, shopGetter: shopGetter}
}

// ensureShopAccess проверяет JWT и что пользователь — владелец магазина (shopId). Возвращает (shopID, nil) или ошибку 401/403.
func (h *Handler) ensureShopAccess(c echo.Context) (int, error) {
	userID := auth_google_handler.UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return 0, c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
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
		return 0, c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}
	return shopID, nil
}

func (h *Handler) getShopIDFromParam(c echo.Context) (int, error) {
	paramID := c.Param("shopId")
	return strconv.Atoi(paramID)
}

func (h *Handler) toResponse(d *domain.ShopAISettings) api.AISettingsResponse {
	if d == nil {
		return api.AISettingsResponse{}
	}
	return api.AISettingsResponse{
		ShopID:      d.ShopID,
		ToneOfVoice: d.ToneOfVoice,
		BrandName:   d.BrandName,
		CustomRules: d.CustomRules,
		IsActive:    d.IsActive,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
