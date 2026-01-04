package notificationendpoint

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) Update(ctx context.Context, d *domain.NotificationEndpoint) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("notification_endpoint update: %w", err)
	}

	d.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *Repository) DisableByID(ctx context.Context, id int) error {
	if err := r.DB.WithContext(ctx).
		Model(&model.ShopNotificationEndpoint{}).
		Where("id = ?", id).
		Update("enabled", false).Error; err != nil {

		return fmt.Errorf("shop_notification_endpoints disable: %w", err)
	}
	return nil
}
