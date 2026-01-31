package saleevent

import (
	"context"
	"fmt"
	"time"

	"salehelp/internal/storage/pggorm/model"
)

// MarkNotified проставляет notified_at для указанных ID событий (в рамках одного shop_id).
func (r *Repository) MarkNotified(ctx context.Context, eventIDs []int, at time.Time) error {
	if len(eventIDs) == 0 {
		return nil
	}
	err := r.DB.WithContext(ctx).
		Model(&model.SaleEvent{}).
		Where("id IN ?", eventIDs).
		Update("notified_at", at).Error
	if err != nil {
		return fmt.Errorf("sale_events mark notified: %w", err)
	}
	return nil
}
