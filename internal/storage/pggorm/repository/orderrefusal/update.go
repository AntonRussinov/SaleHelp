package orderrefusal

import (
	"context"
	"fmt"
	"time"

	"salehelp/internal/storage/pggorm/model"
)

// MarkNotified проставляет notified_at для указанных ID.
func (r *Repository) MarkNotified(ctx context.Context, ids []int, at time.Time) error {
	if len(ids) == 0 {
		return nil
	}
	err := r.DB.WithContext(ctx).
		Model(&model.OrderRefusal{}).
		Where("id IN ?", ids).
		Update("notified_at", at).Error
	if err != nil {
		return fmt.Errorf("order_refusals mark notified: %w", err)
	}
	return nil
}

// MarkAllNotifiedByShop помечает все отказы магазина как уведомлённые (при первом синке).
func (r *Repository) MarkAllNotifiedByShop(ctx context.Context, shopID int, now time.Time) error {
	err := r.DB.WithContext(ctx).
		Model(&model.OrderRefusal{}).
		Where("shop_id = ? AND notified_at IS NULL", shopID).
		Update("notified_at", now).Error
	if err != nil {
		return fmt.Errorf("order_refusals mark all notified by shop: %w", err)
	}
	return nil
}
