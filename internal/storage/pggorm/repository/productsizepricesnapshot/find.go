package productsizepricesnapshot

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.ProductSizePriceSnapshot, error) {
	var m model.ProductSizePriceSnapshot

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("product_size_price_snapshot get by id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListBySnapshotID(ctx context.Context, snapshotID int) ([]*domain.ProductSizePriceSnapshot, error) {
	var rows []model.ProductSizePriceSnapshot

	if err := r.DB.WithContext(ctx).
		Where("product_snapshot_id = ?", snapshotID).
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_size_price_snapshot list by snapshot: %w", err)
	}

	out := make([]*domain.ProductSizePriceSnapshot, 0, len(rows))
	for _, m := range rows {
		out = append(out, m.ToDomain())
	}
	return out, nil
}

func (r *Repository) ListBySnapshotIDs(ctx context.Context, snapshotIDs []int) ([]*domain.ProductSizePriceSnapshot, error) {
	if len(snapshotIDs) == 0 {
		return nil, nil
	}

	var rows []model.ProductSizePriceSnapshot

	if err := r.DB.WithContext(ctx).
		Where("product_snapshot_id IN ?", snapshotIDs).
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_size_price_snapshot list by snapshot_ids: %w", err)
	}

	out := make([]*domain.ProductSizePriceSnapshot, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}
