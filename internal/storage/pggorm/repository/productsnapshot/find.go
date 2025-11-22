package productsnapshot

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) GetLastByProductID(ctx context.Context, productID int) (*domain.ProductSnapshot, error) {
	var m model.ProductSnapshot

	err := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("snapshot_time DESC").
		Limit(1).
		Find(&m).Error

	if err != nil {
		return nil, fmt.Errorf("get last snapshot: %w", err)
	}

	if m.ID == 0 {
		return nil, nil
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListByProductID(ctx context.Context, productID int, limit, offset int) ([]*domain.ProductSnapshot, error) {
	var rows []model.ProductSnapshot

	q := r.DB.WithContext(ctx).Where("product_id = ?", productID).Order("snapshot_time DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list snapshots: %w", err)
	}

	out := make([]*domain.ProductSnapshot, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}
