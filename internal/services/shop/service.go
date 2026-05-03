package shop

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"salehelp/domain"
)

// trialDuration — длительность пробной подписки при регистрации магазина (1.5 месяца = 45 дней).
const trialDuration = 45 * 24 * time.Hour

type ShopsRepo interface {
	Create(ctx context.Context, s *domain.Shop) error
}

type AlertsRepo interface {
	Create(ctx context.Context, a *domain.Alert) error
}

type SubscriptionsRepo interface {
	Create(ctx context.Context, s *domain.Subscription) error
}

type Service struct {
	shops  ShopsRepo
	alerts AlertsRepo
	subs   SubscriptionsRepo
	log    *slog.Logger
}

func New(shops ShopsRepo, alerts AlertsRepo, subs SubscriptionsRepo, log *slog.Logger) *Service {
	return &Service{
		shops:  shops,
		alerts: alerts,
		subs:   subs,
		log:    log.With("service", "shop"),
	}
}

// CreateShop создаёт магазин и автоматически настраивает стандартные алерты и пробную подписку.
func (s *Service) CreateShop(ctx context.Context, shop *domain.Shop) error {
	if err := s.shops.Create(ctx, shop); err != nil {
		return err
	}
	s.setupDefaultAlerts(ctx, shop.ID)
	s.createTrialSubscription(ctx, shop.ID)
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

func (s *Service) createTrialSubscription(ctx context.Context, shopID int) {
	now := time.Now()
	sub := &domain.Subscription{
		ShopID:          shopID,
		Plan:            "basic",
		Status:          "active",
		ValidFrom:       now,
		ValidTo:         now.Add(trialDuration),
		BaseAmountMinor: 0,
		CurrencyCode:    "RUB",
		Source:          "auto_trial",
	}
	if err := s.subs.Create(ctx, sub); err != nil {
		s.log.Error("failed to create trial subscription", "shopID", shopID, "error", err)
	}
}
