package shop

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"salehelp/domain"
)

const trialDays = 14

// defaultFreeShopLimit — лимит магазинов без активной подписки (после истечения триала).
const defaultFreeShopLimit = 1

type ShopsRepo interface {
	Create(ctx context.Context, s *domain.Shop) error
	CountByUser(ctx context.Context, userID int) (int, error)
}

type AlertsRepo interface {
	Create(ctx context.Context, a *domain.Alert) error
}

type SubscriptionsRepo interface {
	GetActiveByUserID(ctx context.Context, userID int) (*domain.Subscription, error)
	Create(ctx context.Context, s *domain.Subscription) error
}

type UsersRepo interface {
	GetByID(ctx context.Context, id int) (*domain.User, error)
}

type EndpointsRepo interface {
	Exists(ctx context.Context, shopID int, channel, address string) (bool, error)
	Create(ctx context.Context, e *domain.NotificationEndpoint) error
}

type Service struct {
	shops     ShopsRepo
	alerts    AlertsRepo
	subs      SubscriptionsRepo
	users     UsersRepo
	endpoints EndpointsRepo
	log       *slog.Logger
}

func New(shops ShopsRepo, alerts AlertsRepo, subs SubscriptionsRepo, log *slog.Logger) *Service {
	return &Service{
		shops:  shops,
		alerts: alerts,
		subs:   subs,
		log:    log.With("service", "shop"),
	}
}

// WithUserEndpoints wires optional repos for auto-creating Telegram notification endpoints.
func (s *Service) WithUserEndpoints(users UsersRepo, endpoints EndpointsRepo) *Service {
	s.users = users
	s.endpoints = endpoints
	return s
}

// CreateShop создаёт магазин, стандартные алерты и пробную подписку (если у пользователя ещё нет активной).
func (s *Service) CreateShop(ctx context.Context, shop *domain.Shop) error {
	if err := s.checkShopLimit(ctx, shop.UserID); err != nil {
		return err
	}
	if err := s.shops.Create(ctx, shop); err != nil {
		return err
	}
	s.setupDefaultAlerts(ctx, shop.ID)
	s.createTrialIfNeeded(ctx, shop.UserID)
	s.autoCreateTelegramEndpoint(ctx, shop)
	return nil
}

// autoCreateTelegramEndpoint создаёт Telegram-эндпоинт для нового магазина,
// если у владельца уже привязан Telegram-чат (через deep link или бот).
func (s *Service) autoCreateTelegramEndpoint(ctx context.Context, shop *domain.Shop) {
	if s.users == nil || s.endpoints == nil {
		return
	}
	user, err := s.users.GetByID(ctx, shop.UserID)
	if err != nil || user == nil || user.TelegramChatID == nil {
		return
	}
	chatStr := fmt.Sprintf("%d", *user.TelegramChatID)
	exists, err := s.endpoints.Exists(ctx, shop.ID, "telegram", chatStr)
	if err != nil || exists {
		return
	}
	label := "owner"
	ep := &domain.NotificationEndpoint{
		ShopID:  shop.ID,
		Channel: "telegram",
		Address: chatStr,
		Enabled: true,
		Label:   &label,
	}
	if err := s.endpoints.Create(ctx, ep); err != nil {
		s.log.Error("auto-create telegram endpoint failed", "shop_id", shop.ID, "user_id", shop.UserID, "err", err)
	}
}

// checkShopLimit проверяет, не превышен ли лимит магазинов по подписке.
func (s *Service) checkShopLimit(ctx context.Context, userID int) error {
	count, err := s.shops.CountByUser(ctx, userID)
	if err != nil {
		return err
	}

	sub, err := s.subs.GetActiveByUserID(ctx, userID)
	if err != nil {
		return err
	}

	maxShops := defaultFreeShopLimit
	if sub != nil {
		maxShops = sub.MaxShops + sub.ExtraShops
	}

	if count >= maxShops {
		return fmt.Errorf("достигнут лимит магазинов: ваш тариф позволяет %d магазин(ов), у вас уже %d", maxShops, count)
	}
	return nil
}

func (s *Service) setupDefaultAlerts(ctx context.Context, shopID int) {
	stockLowPayload, _ := json.Marshal(map[string]any{"threshold": 5})
	priceDropPayload, _ := json.Marshal(map[string]any{"percent": 3.0, "mode": "any"})

	type alertDef struct {
		alertType string
		severity  int
		message   string
		payload   []byte
	}

	defaults := []alertDef{
		{domain.AlertTypeStockZero, 5, "Товар закончился", nil},
		{domain.AlertTypeStockLow, 2, "Мало остатков (< 5 шт.)", stockLowPayload},
		{domain.AlertTypePriceDrop, 2, "Снижение цены > 3%", priceDropPayload},
	}

	for _, def := range defaults {
		alert := &domain.Alert{
			ShopID:   shopID,
			Type:     def.alertType,
			Severity: def.severity,
			Message:  def.message,
			Payload:  def.payload,
			Scope:    domain.AlertScopeAuto,
			Status:   domain.AlertStatusActive,
		}
		if err := s.alerts.Create(ctx, alert); err != nil {
			s.log.Error("failed to create default alert", "shopID", shopID, "type", def.alertType, "error", err)
		}
	}
}

// createTrialIfNeeded выдаёт 14-дневный триал если у пользователя ещё нет активной подписки.
func (s *Service) createTrialIfNeeded(ctx context.Context, userID int) {
	existing, err := s.subs.GetActiveByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to check active subscription", "userID", userID, "error", err)
		return
	}
	if existing != nil {
		return
	}

	cfg := domain.TierConfigs[domain.TierTrial]
	now := time.Now()
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	sub := &domain.Subscription{
		UserID:       userID,
		Tier:         domain.TierTrial,
		Status:       domain.SubStatusTrial,
		PeriodMonths: 0,
		StartedAt:    now,
		ExpiresAt:    now.AddDate(0, 0, trialDays),
		AIQuota:      cfg.AIQuotaMonthly,
		AIUsed:       0,
		AIResetsAt:   nextMonth,
		MaxShops:     cfg.MaxShops,
		ExtraShops:   0,
		PriceMinor:   0,
	}
	if err := s.subs.Create(ctx, sub); err != nil {
		s.log.Error("failed to create trial subscription", "userID", userID, "error", err)
	}
}
