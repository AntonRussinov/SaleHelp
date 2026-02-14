package shop_ai_settings

import (
	"context"
	"salehelp/api"
	"salehelp/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

// AISettingsRepo интерфейс репозитория настроек AI.
type AISettingsRepo interface {
	GetByShopID(ctx context.Context, shopID int) (*domain.ShopAISettings, error)
	Save(ctx context.Context, d *domain.ShopAISettings) error
}

// Handler HTTP-хендлеры для настроек AI магазина (бренд, тон, правила для промпта).
type Handler struct {
	repo AISettingsRepo
}

// NewHandler создаёт хендлер.
func NewHandler(repo AISettingsRepo) *Handler {
	return &Handler{repo: repo}
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
