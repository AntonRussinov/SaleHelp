package wbwarehouse

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByWBID(ctx context.Context, wbID int64) (*domain.WBWarehouse, error) {
	var m model.WBWarehouse

	if err := r.DB.WithContext(ctx).
		First(&m, "wb_id = ?", wbID).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("wb_warehouse get by wb_id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListAll(ctx context.Context) ([]*domain.WBWarehouse, error) {
	var rows []model.WBWarehouse

	if err := r.DB.WithContext(ctx).
		Order("id").
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("wb_warehouse list: %w", err)
	}

	result := make([]*domain.WBWarehouse, 0, len(rows))
	for _, r := range rows {
		result = append(result, r.ToDomain())
	}

	return result, nil
}
