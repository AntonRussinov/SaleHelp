package pricerestorebatch

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) SetRestoredAt(ctx context.Context, batchID int, at time.Time) error {
	result := r.DB.WithContext(ctx).
		Table("price_restore_batches").
		Where("id = ?", batchID).
		Update("restored_at", at)
	if result.Error != nil {
		return fmt.Errorf("price_restore_batch set restored_at: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("price_restore_batch set restored_at: batch id=%d not found", batchID)
	}
	return nil
}
