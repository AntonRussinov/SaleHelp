package ai_settings

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"salehelp/api"
	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"

	"github.com/labstack/echo/v4"
)

type AISettingsRepo interface {
	GetByShopID(ctx context.Context, shopID int) (*domain.AISettings, error)
	Upsert(ctx context.Context, d *domain.AISettings) error
	ListRules(ctx context.Context, settingsID int) ([]*domain.AIRule, error)
	CreateRule(ctx context.Context, d *domain.AIRule) error
	UpdateRule(ctx context.Context, d *domain.AIRule) error
	DeleteRule(ctx context.Context, settingsID, ruleID int) error
	ListTemplates(ctx context.Context, settingsID int) ([]*domain.AITemplate, error)
	UpsertTemplate(ctx context.Context, d *domain.AITemplate) error
}

type ShopGetter interface {
	GetByID(ctx context.Context, id int) (*domain.Shop, error)
}

type SubscriptionsRepo interface {
	GetActiveByUserID(ctx context.Context, userID int) (*domain.Subscription, error)
}

type Handler struct {
	repo     AISettingsRepo
	shopRepo ShopGetter
	subsRepo SubscriptionsRepo
}

func NewHandler(repo AISettingsRepo, shopRepo ShopGetter, subsRepo SubscriptionsRepo) *Handler {
	return &Handler{repo: repo, shopRepo: shopRepo, subsRepo: subsRepo}
}

func (h *Handler) ensureShopAccess(c echo.Context) (int, *domain.Shop, error) {
	userID := auth_google_handler.UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return 0, nil, c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	shopID, err := strconv.Atoi(c.Param("shopId"))
	if err != nil || shopID <= 0 {
		return 0, nil, echo.NewHTTPError(http.StatusBadRequest, "invalid shop ID")
	}
	shop, err := h.shopRepo.GetByID(c.Request().Context(), shopID)
	if err != nil {
		return 0, nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if shop == nil || shop.UserID != userID {
		return 0, nil, c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}
	return shopID, shop, nil
}

// maxRulesForTier возвращает лимит правил по тарифу (0 = запрещено).
func maxRulesForTier(tier string) int {
	cfg := domain.TierConfigs[tier]
	return cfg.MaxAIRules
}

func settingsToResponse(s *domain.AISettings, rules []*domain.AIRule, templates []*domain.AITemplate) api.AISettingsResponse {
	ruleResp := make([]api.AIRuleResponse, len(rules))
	for i, r := range rules {
		ruleResp[i] = api.AIRuleResponse{
			ID:         r.ID,
			RuleType:   r.RuleType,
			RuleText:   r.RuleText,
			OrderIndex: r.OrderIndex,
		}
	}
	tmplResp := make([]api.AITemplateResponse, len(templates))
	for i, t := range templates {
		tmplResp[i] = api.AITemplateResponse{
			ID:                t.ID,
			Category:          t.Category,
			CustomPrompt:      t.CustomPrompt,
			ToneVoiceOverride: t.ToneVoiceOverride,
			Enabled:           t.Enabled,
		}
	}
	return api.AISettingsResponse{
		ID:                 s.ID,
		ShopID:             s.ShopID,
		IsActive:           s.IsActive,
		ToneOfVoice:        s.ToneOfVoice,
		ResponseLength:     s.ResponseLength,
		Addressing:         s.Addressing,
		EmojiMode:          s.EmojiMode,
		BrandName:          s.BrandName,
		BrandDescription:   s.BrandDescription,
		RepresentativeName: s.RepresentativeName,
		Signature:          s.Signature,
		Rules:              ruleResp,
		Templates:          tmplResp,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
}

func defaultSettings(shopID, userID int) *domain.AISettings {
	return &domain.AISettings{
		ShopID:         shopID,
		UserID:         userID,
		IsActive:       false,
		ToneOfVoice:    "friendly",
		ResponseLength: "medium",
		Addressing:     "vy",
		EmojiMode:      "minimal",
		CreatedAt:      time.Now(),
	}
}
