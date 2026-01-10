package subscriptionaddon

import (
	"context"
	"fmt"
	"salehelp/domain"

	"gorm.io/gorm/clause"
)

// Upsert: 1 строка на (subscription_id, code)
func (r *Repository) Upsert(ctx context.Context, a *domain.SubscriptionAddon) error {
	m := toModel(a)

	err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "subscription_id"},
				{Name: "code"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"quantity",
				"amount_minor",
				"currency_code",
				"meta",
			}),
		}).
		Create(m).Error

	if err != nil {
		return fmt.Errorf("subscription addon upsert: %w", err)
	}

	a.ID = m.ID
	a.CreatedAt = m.CreatedAt
	return nil
}
