package pricerestorebatch

import (
	"context"
	"fmt"
	"salehelp/domain"
	"time"
)

func (r *Repository) SetRestoredAt(ctx context.Context, batchID int, at time.Time) error {
	result := r.DB.WithContext(ctx).
		Table("price_restore_batches").
		Where("id = ?", batchID).
		Updates(map[string]any{
			"restored_at":    at,
			"restore_status": domain.RestoreStatusRestored,
		})
	if result.Error != nil {
		return fmt.Errorf("price_restore_batch set restored_at: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("price_restore_batch set restored_at: batch id=%d not found", batchID)
	}
	return nil
}

// TrySetRestoring атомарно переводит batch из 'pending' в 'restoring'.
// Возвращает true если успешно (никто другой не начал восстановление).
func (r *Repository) TrySetRestoring(ctx context.Context, batchID int) (bool, error) {
	now := time.Now()
	result := r.DB.WithContext(ctx).
		Table("price_restore_batches").
		Where("id = ? AND restore_status = ?", batchID, domain.RestoreStatusPending).
		Updates(map[string]any{
			"restore_status":    domain.RestoreStatusRestoring,
			"restore_started_at": now,
		})
	if result.Error != nil {
		return false, fmt.Errorf("price_restore_batch try set restoring: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

// SetRestoreFailed помечает batch как failed.
func (r *Repository) SetRestoreFailed(ctx context.Context, batchID int) error {
	result := r.DB.WithContext(ctx).
		Table("price_restore_batches").
		Where("id = ?", batchID).
		Update("restore_status", domain.RestoreStatusFailed)
	if result.Error != nil {
		return fmt.Errorf("price_restore_batch set failed: %w", result.Error)
	}
	return nil
}
