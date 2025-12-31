package notificationendpoint

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Update(ctx context.Context, d *domain.NotificationEndpoint) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("notification_endpoint update: %w", err)
	}

	d.UpdatedAt = m.UpdatedAt
	return nil
}
