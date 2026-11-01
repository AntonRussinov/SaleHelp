package pricerestorebatch

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.PriceRestoreBatch, error) {
	var m model.PriceRestoreBatch
	if err := r.DB.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("price_restore_batch get by id: %w", err)
	}
	return m.ToDomain(), nil
}

// GetTriggeredIDsByBatchID возвращает список triggered_id для пакета.
func (r *Repository) GetTriggeredIDsByBatchID(ctx context.Context, batchID int) ([]int, error) {
	var ids []int
	err := r.DB.WithContext(ctx).
		Table("price_restore_batch_triggered").
		Where("batch_id = ?", batchID).
		Pluck("triggered_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("price_restore_batch_triggered get by batch: %w", err)
	}
	return ids, nil
}
