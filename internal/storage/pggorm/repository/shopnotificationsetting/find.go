package shopnotificationsetting

import (
	"context"
	"fmt"

	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) ListByShopID(ctx context.Context, shopID int) ([]*model.ShopNotificationSetting, error) {
	if shopID <= 0 {
		return nil, fmt.Errorf("shop_notification_setting list: invalid shop_id=%d", shopID)
	}

	var rows []*model.ShopNotificationSetting
	if err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("shop_notification_setting list: %w", err)
	}
	return rows, nil
}

func (r *Repository) MapByShopID(ctx context.Context, shopID int) (map[string]bool, error) {
	rows, err := r.ListByShopID(ctx, shopID)
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(rows))
	for _, row := range rows {
		if row == nil || row.EventKey == "" {
			continue
		}
		out[row.EventKey] = row.Enabled
	}
	return out, nil
}
