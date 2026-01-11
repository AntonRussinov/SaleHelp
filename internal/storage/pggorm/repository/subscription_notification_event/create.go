package subscriptionnotificationevent

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

func (r *Repository) CreateOnce(ctx context.Context, subscriptionID int, event string) (bool, error) {
	e := &model.SubscriptionNotificationEvent{
		SubscriptionID: subscriptionID,
		Event:          event,
	}

	// ON CONFLICT DO NOTHING
	tx := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "subscription_id"}, {Name: "event"}},
			DoNothing: true,
		}).
		Create(e)

	if tx.Error != nil {
		return false, fmt.Errorf("create subscription notification event: %w", tx.Error)
	}

	// RowsAffected = 1 если вставили, 0 если конфликт
	return tx.RowsAffected == 1, nil
}
