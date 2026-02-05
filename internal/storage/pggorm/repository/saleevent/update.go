package saleevent

import (
	"context"
	"fmt"
	"time"

	"salehelp/internal/storage/pggorm/model"
)

// MarkNotified проставляет notified_at для указанных ID событий.
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

// MarkAllNotifiedByShop помечает все sale_events магазина как уведомлённые (notified_at = now).
// Используется при первом синке, чтобы не отправлять историю выкупов/возвратов.
func (r *Repository) MarkAllNotifiedByShop(ctx context.Context, shopID int, now time.Time) error {
	err := r.DB.WithContext(ctx).
		Model(&model.SaleEvent{}).
		Where("shop_id = ? AND notified_at IS NULL", shopID).
		Update("notified_at", now).Error
	if err != nil {
		return fmt.Errorf("sale_events mark all notified by shop: %w", err)
	}
	return nil
}
