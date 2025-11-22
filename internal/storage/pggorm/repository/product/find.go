package product

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	var m model.Product

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("product get by id: %w", err)
	}

	return m.ToDomain(), nil
}

// GetByShopAndExternalID получает продукт по (shop_id, external_id)
// это твой уникальный ключ
func (r *Repository) GetByShopAndExternalID(
	ctx context.Context,
	shopID int,
	externalID string,
) (*domain.Product, error) {
	var m model.Product

	if err := r.DB.WithContext(ctx).
		First(&m, "shop_id = ? AND external_id = ?", shopID, externalID).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("product get by shop_id+external_id: %w", err)
	}

	return m.ToDomain(), nil
}

// ListByShop список продуктов конкретного магазина c пагинацией
func (r *Repository) ListByShop(
	ctx context.Context,
	shopID int,
	limit, offset int,
) ([]*domain.Product, error) {
	var models []model.Product

	q := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Order("id ASC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("product list by shop: %w", err)
	}

	result := make([]*domain.Product, 0, len(models))
	for i := range models {
		result = append(result, models[i].ToDomain())
	}

	return result, nil
}
