package notificationendpoint

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) CountEnabledByShopAndChannel(ctx context.Context, shopID int, channel string) (int, error) {
	if shopID <= 0 {
		return 0, fmt.Errorf("shop_id is required")
	}
	if channel == "" {
		return 0, fmt.Errorf("channel is required")
	}

	var cnt int64
	err := r.DB.WithContext(ctx).
		Model(&model.ShopNotificationEndpoint{}).
		Where("shop_id = ? AND channel = ? AND enabled = TRUE", shopID, channel).
		Count(&cnt).Error
	if err != nil {
		return 0, fmt.Errorf("count enabled endpoints by shop/channel: %w", err)
	}

	return int(cnt), nil
}
