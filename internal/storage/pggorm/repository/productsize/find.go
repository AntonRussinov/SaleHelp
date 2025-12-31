package productsize

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.ProductSize, error) {
	var m model.ProductSize

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("product_size get by id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListByProductID(ctx context.Context, productID int) ([]*domain.ProductSize, error) {
	var rows []model.ProductSize

	if err := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_size list by product: %w", err)
	}

	out := make([]*domain.ProductSize, 0, len(rows))
	for _, m := range rows {
		out = append(out, m.ToDomain())
	}
	return out, nil
}

func (r *Repository) ListByProductIDs(
	ctx context.Context,
	productIDs []int64,
) ([]*domain.ProductSize, error) {

	if len(productIDs) == 0 {
		return []*domain.ProductSize{}, nil
	}

	var rows []model.ProductSize

	if err := r.DB.WithContext(ctx).
		Where("product_id IN ?", productIDs).
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_size list by product ids: %w", err)
	}

	out := make([]*domain.ProductSize, 0, len(rows))
	for _, m := range rows {
		out = append(out, m.ToDomain())
	}

	return out, nil
}

func (r *Repository) FindByProductAndExternalID(ctx context.Context, productID int, externalID string) (*domain.ProductSize, error) {
	var m model.ProductSize

	if err := r.DB.WithContext(ctx).
		Where("product_id = ? AND external_id = ?", productID, externalID).
		First(&m).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("product_size find by product+external: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListByProductSizeID(ctx context.Context, productSizeID int) ([]*domain.ProductSizeSKU, error) {
	var rows []model.ProductSizeSKU

	if err := r.DB.WithContext(ctx).
		Where("product_size_id = ?", productSizeID).
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_size_sku list by size: %w", err)
	}

	out := make([]*domain.ProductSizeSKU, 0, len(rows))
	for _, m := range rows {
		out = append(out, m.ToDomain())
	}
	return out, nil
}
