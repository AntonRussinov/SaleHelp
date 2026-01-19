package notificationendpoint

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Create(ctx context.Context, d *domain.NotificationEndpoint) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("notification_endpoint create: %w", err)
	}

	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt
	return nil
}
