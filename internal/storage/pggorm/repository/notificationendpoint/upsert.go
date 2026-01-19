package notificationendpoint

import (
	"context"
	"fmt"
	"salehelp/domain"

	"gorm.io/gorm/clause"
)

func (r *Repository) Upsert(ctx context.Context, d *domain.NotificationEndpoint) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "shop_id"},
				{Name: "channel"},
				{Name: "address"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"enabled",
				"label",
				"updated_at",
			}),
		}).
		Create(m).Error; err != nil {
		return fmt.Errorf("notification_endpoint upsert: %w", err)
	}

	// на insert вернёт ID, на update обычно нет, но не мешает
	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt
	return nil
}
