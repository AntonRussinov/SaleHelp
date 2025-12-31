package alert

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.Alert, error) {
	var m model.Alert

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("alert get by id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListByShop(ctx context.Context, shopID int, limit, offset int) ([]*domain.Alert, error) {
	var rows []model.Alert

	q := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Order("triggered_at DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("alert list by shop: %w", err)
	}

	out := make([]*domain.Alert, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}

func (r *Repository) ListByProduct(ctx context.Context, productID int, limit, offset int) ([]*domain.Alert, error) {
	var rows []model.Alert

	q := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("triggered_at DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("alert list by product: %w", err)
	}

	out := make([]*domain.Alert, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}

func (r *Repository) ListActiveByShop(ctx context.Context, shopID int) ([]*domain.Alert, error) {
	var rows []model.Alert

	if err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Where("status IN ?", []string{"active", "new"}).
		Order("id ASC").
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("alert list active by shop: %w", err)
	}

	out := make([]*domain.Alert, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain()) // если ToDomain() возвращает *domain.Alert
	}
	return out, nil
}

func (r *Repository) ListActiveByShopAndType(ctx context.Context, shopID int, alertType string) ([]*domain.Alert, error) {
	var rows []model.Alert

	if err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Where("type = ?", alertType).
		Where("status IN ?", []string{"active", "new"}).
		Order("id ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("alert list active by shop and type: %w", err)
	}

	out := make([]*domain.Alert, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}
