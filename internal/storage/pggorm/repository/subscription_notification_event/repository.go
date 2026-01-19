package subscriptionnotificationevent

import (
	"log/slog"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm"
	"salehelp/internal/storage/pggorm/model"
	"salehelp/pkg/logger/sl"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	pggorm.Repository
}

func NewRepository(log *slog.Logger, db *gorm.DB) *Repository {
	repo := &Repository{}
	repo.Log = log.With(sl.Module("SubscriptionNotificationEventsRepository"))
	repo.DB = db
	return repo
}

func (r *Repository) WithDB(db *gorm.DB) *Repository {
	cp := *r
	cp.DB = db
	return &cp
}

// domain → model
func toModel(e *domain.SubscriptionNotificationEvent) *model.SubscriptionNotificationEvent {
	if e == nil {
		return nil
	}
	return &model.SubscriptionNotificationEvent{
		ID:             e.ID,
		SubscriptionID: e.SubscriptionID,
		Event:          e.Event,
		SentAt:         e.SentAt,
	}
}

// (если понадобится, но обычно ты пользуешься e.ToDomain() у модели)
func toDomain(m *model.SubscriptionNotificationEvent) *domain.SubscriptionNotificationEvent {
	if m == nil {
		return nil
	}
	return &domain.SubscriptionNotificationEvent{
		ID:             m.ID,
		SubscriptionID: m.SubscriptionID,
		Event:          m.Event,
		SentAt:         m.SentAt,
	}
}

// Для единообразия: сейчас не используется, но пусть будет.
// Удобно в тестах (args) или если добавишь Create() позже.
func saveSQLArgs(e domain.SubscriptionNotificationEvent) []any {
	return []any{
		e.SubscriptionID,
		e.Event,
		// sent_at autoCreateTime, но если в тестах ожидаешь явно — можно передавать:
		// e.SentAt,
		time.Time{}, // placeholder если нужно; чаще не надо, см. тесты
	}
}
