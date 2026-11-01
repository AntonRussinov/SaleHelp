package dbworder

import (
	"context"
	"fmt"
	"time"

	"salehelp/internal/storage/pggorm/model"
)

// MarkNotified проставляет notified_at для указанных ID заказов DBW.
func (r *Repository) MarkNotified(ctx context.Context, orderIDs []int, at time.Time) error {
	if len(orderIDs) == 0 {
		return nil
	}
	err := r.DB.WithContext(ctx).
		Model(&model.DBWOrder{}).
		Where("id IN ?", orderIDs).
		Update("notified_at", at).Error
	if err != nil {
		return fmt.Errorf("dbw_orders mark notified: %w", err)
	}
	return nil
}
