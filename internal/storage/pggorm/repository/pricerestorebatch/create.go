package pricerestorebatch

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) Create(ctx context.Context, batch *domain.PriceRestoreBatch, triggeredIDs []int) (int, error) {
	if batch == nil {
		return 0, fmt.Errorf("price_restore_batch create: batch is nil")
	}
	m := &model.PriceRestoreBatch{
		ShopID:        batch.ShopID,
		RestoreStatus: domain.RestoreStatusPending,
		CreatedAt:     batch.CreatedAt,
		RestoredAt:    batch.RestoredAt,
	}
	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return 0, fmt.Errorf("price_restore_batch create: %w", err)
	}
	batch.ID = m.ID

	if len(triggeredIDs) > 0 {
		for _, tid := range triggeredIDs {
			if err := r.DB.WithContext(ctx).Exec(
				"INSERT INTO price_restore_batch_triggered (batch_id, triggered_id) VALUES (?, ?)",
				m.ID, tid,
			).Error; err != nil {
				return 0, fmt.Errorf("price_restore_batch_triggered create: %w", err)
			}
		}
	}

	return m.ID, nil
}
