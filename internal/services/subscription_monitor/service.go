package subscriptionmonitor

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/services/notification_settings"
	"time"
)

const (
	EventExpiring3d = "expiring_3d"
	EventExpiring0d = "expiring_0d"
	EventExpired1d  = "expired_1d"
)

type SubscriptionsRepo interface {
	ListExpiringBetween(ctx context.Context, from, to time.Time) ([]*domain.Subscription, error)
}

type ShopsRepo interface {
	ListByUser(ctx context.Context, userID int, limit, offset int) ([]*domain.Shop, error)
}

type EndpointsRepo interface {
	ListEnabledTelegramChatIDsByShopID(ctx context.Context, shopID int) ([]int64, error)
}

type TelegramAPI interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type Monitor struct {
	subsRepo  SubscriptionsRepo
	shopsRepo ShopsRepo
	endpoints EndpointsRepo
	tg        TelegramAPI
	settings  *notification_settings.Service
	now       func() time.Time
	loc       *time.Location
}

func New(subsRepo SubscriptionsRepo, shopsRepo ShopsRepo, endpoints EndpointsRepo, tg TelegramAPI, settings *notification_settings.Service, loc *time.Location) *Monitor {
	if loc == nil {
		loc = time.Local
	}
	return &Monitor{
		subsRepo:  subsRepo,
		shopsRepo: shopsRepo,
		endpoints: endpoints,
		tg:        tg,
		settings:  settings,
		now:       time.Now,
		loc:       loc,
	}
}

func (m *Monitor) SetNow(fn func() time.Time) { m.now = fn }

func (m *Monitor) Run(ctx context.Context) error {
	now := m.now().In(m.loc)

	dayStart := func(t time.Time) time.Time {
		y, mo, d := t.Date()
		return time.Date(y, mo, d, 0, 0, 0, 0, t.Location())
	}

	startToday := dayStart(now)
	if err := m.process(ctx, startToday.AddDate(0, 0, 3), startToday.AddDate(0, 0, 4), EventExpiring3d); err != nil {
		return err
	}
	if err := m.process(ctx, startToday, startToday.AddDate(0, 0, 1), EventExpiring0d); err != nil {
		return err
	}
	if err := m.process(ctx, startToday.AddDate(0, 0, -1), startToday, EventExpired1d); err != nil {
		return err
	}
	return nil
}

func (m *Monitor) process(ctx context.Context, from, to time.Time, event string) error {
	subs, err := m.subsRepo.ListExpiringBetween(ctx, from, to)
	if err != nil {
		return fmt.Errorf("list subs for %s: %w", event, err)
	}

	for _, s := range subs {
		shops, err := m.shopsRepo.ListByUser(ctx, s.UserID, 20, 0)
		if err != nil {
			return fmt.Errorf("list shops for user %d: %w", s.UserID, err)
		}

		notified := map[int64]struct{}{}
		for _, shop := range shops {
			if m.settings != nil {
				enabled, err := m.settings.EffectiveEnabled(ctx, shop.ID, domain.NotificationEventSubscriptionStatus)
				if err == nil && !enabled {
					continue
				}
			}

			chatIDs, err := m.endpoints.ListEnabledTelegramChatIDsByShopID(ctx, shop.ID)
			if err != nil {
				return fmt.Errorf("list chat ids for shop %d: %w", shop.ID, err)
			}

			text := buildText(event, s.ExpiresAt.In(m.loc))
			for _, chatID := range chatIDs {
				if _, already := notified[chatID]; already {
					continue
				}
				notified[chatID] = struct{}{}
				_ = m.tg.SendMessage(ctx, chatID, text)
			}
		}
	}
	return nil
}

func buildText(event string, expiresAt time.Time) string {
	d := expiresAt.Format("02.01.2006")
	switch event {
	case EventExpiring3d:
		return fmt.Sprintf("⏳ Подписка истекает через 3 дня (%s).\nЕсли не продлить — алерты остановятся.", d)
	case EventExpiring0d:
		return fmt.Sprintf("⚠️ Сегодня истекает подписка (%s).\nЕсли не продлить — алерты будут остановлены.", d)
	case EventExpired1d:
		return fmt.Sprintf("⛔ Подписка истекла (%s).\nАлерты остановлены до продления подписки.", d)
	default:
		return "Уведомление по подписке"
	}
}
