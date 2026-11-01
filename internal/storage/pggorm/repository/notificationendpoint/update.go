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

// UpdateAddressByUser меняет address (chat_id) в эндпоинтах пользователя.
// Используется при перепривязке Telegram: все эндпоинты со старым chat_id обновляются на новый.
func (r *Repository) UpdateAddressByUser(ctx context.Context, userID int, channel, oldAddress, newAddress string) error {
	err := r.DB.WithContext(ctx).Exec(`
		UPDATE shop_notification_endpoints
		SET address = ?
		WHERE channel = ? AND address = ?
		  AND shop_id IN (SELECT id FROM shops WHERE user_id = ?)`,
		newAddress, channel, oldAddress, userID,
	).Error
	if err != nil {
		return fmt.Errorf("endpoint update address by user: %w", err)
	}
	return nil
}

func (r *Repository) SetEnabledByShopAndAddress(
	ctx context.Context,
	shopID int,
	channel string,
	address string,
	enabled bool,
) error {

	err := r.DB.WithContext(ctx).
		Model(&model.ShopNotificationEndpoint{}).
		Where("shop_id = ? AND channel = ? AND address = ?", shopID, channel, address).
		Update("enabled", enabled).Error

	if err != nil {
		return fmt.Errorf("endpoint set enabled: %w", err)
	}
	return nil
}
