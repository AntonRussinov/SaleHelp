package subscriptionmonitor

import (
	"context"
	"fmt"
	"salehelp/domain"
	"time"
)

const (
	EventExpiring3d = "expiring_3d"
	EventExpiring0d = "expiring_0d"
	EventExpired1d  = "expired_1d"
)

type SubscriptionsRepo interface {
	ListByValidToRange(ctx context.Context, from, to time.Time) ([]*domain.Subscription, error)
}

type SubsriptionNotificationEventsRepo interface {
	CreateOnce(ctx context.Context, subscriptionID int, event string) (inserted bool, err error)
}

type EndpointsRepo interface {
	ListEnabledTelegramChatIDsByShopID(ctx context.Context, shopID int) ([]int64, error)
}

type TelegramAPI interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type Monitor struct {
	subsRepo   SubscriptionsRepo
	eventsRepo SubsriptionNotificationEventsRepo
	endpoints  EndpointsRepo
	tg         TelegramAPI
	now        func() time.Time
	loc        *time.Location
}

func New(subsRepo SubscriptionsRepo, eventsRepo SubsriptionNotificationEventsRepo, endpoints EndpointsRepo, tg TelegramAPI, loc *time.Location) *Monitor {
	if loc == nil {
		loc = time.Local
	}
	return &Monitor{
		subsRepo: subsRepo, eventsRepo: eventsRepo, endpoints: endpoints, tg: tg,
		now: time.Now,
		loc: loc,
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
	startTomorrow := startToday.AddDate(0, 0, 1)

	startIn3 := startToday.AddDate(0, 0, 3)
	startIn4 := startToday.AddDate(0, 0, 4)

	startYesterday := startToday.AddDate(0, 0, -1)

	if err := m.process(ctx, startIn3, startIn4, EventExpiring3d); err != nil {
		return err
	}
	if err := m.process(ctx, startToday, startTomorrow, EventExpiring0d); err != nil {
		return err
	}
	if err := m.process(ctx, startYesterday, startToday, EventExpired1d); err != nil {
		return err
	}

	return nil
}

func (m *Monitor) process(ctx context.Context, from, to time.Time, event string) error {
	subs, err := m.subsRepo.ListByValidToRange(ctx, from, to)
	if err != nil {
		return fmt.Errorf("list subs for %s: %w", event, err)
	}

	for _, s := range subs {
		inserted, err := m.eventsRepo.CreateOnce(ctx, s.ID, event)
		if err != nil {
			return fmt.Errorf("create event %s: %w", event, err)
		}
		if !inserted {
			continue
		}

		chatIDs, err := m.endpoints.ListEnabledTelegramChatIDsByShopID(ctx, s.ShopID)
		if err != nil {
			return fmt.Errorf("list chat ids: %w", err)
		}
		if len(chatIDs) == 0 {
			continue
		}

		text := buildText(event, s.ValidTo.In(m.loc))

		for _, chatID := range chatIDs {
			_ = m.tg.SendMessage(ctx, chatID, text) // v0: best-effort
		}
	}

	return nil
}

func buildText(event string, validTo time.Time) string {
	d := validTo.Format("02.01.2006")
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
