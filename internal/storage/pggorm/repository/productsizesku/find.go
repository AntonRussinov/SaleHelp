package productsizesku

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

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

// func (r *Repository) ListByProductSizeIDs(ctx context.Context, productSizeIDs []int) ([]*domain.ProductSizeSKU, error) {
// 	var rows []model.ProductSizeSKU

// 	if err := r.DB.WithContext(ctx).
// 		Where("product_size_id IN ?", productSizeIDs).
// 		Find(&rows).Error; err != nil {

// 		return nil, fmt.Errorf("product_size_sku list by size: %w", err)
// 	}

// 	out := make([]*domain.ProductSizeSKU, 0, len(rows))
// 	for _, m := range rows {
// 		out = append(out, m.ToDomain())
// 	}
// 	return out, nil
// }

func (r *Repository) ListByProductSizeIDs(ctx context.Context, productSizeIDs []int) (map[int][]string, error) {
	if len(productSizeIDs) == 0 {
		return map[int][]string{}, nil
	}

	var rows []model.ProductSizeSKU

	if err := r.DB.WithContext(ctx).
		Where("product_size_id IN ?", productSizeIDs).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("product_size_sku list by size ids: %w", err)
	}

	out := make(map[int][]string, len(productSizeIDs))
	for _, r := range rows {
		out[r.ProductSizeID] = append(out[r.ProductSizeID], r.SKU)
	}

	return out, nil
}

// Stocks sync helper: SKU + ссылки
type SKUWithRefs struct {
	SKU           string
	ProductID     int
	ProductSizeID int
}

func (r *Repository) ListByProductIDs(ctx context.Context, productIDs []int64) ([]SKUWithRefs, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}

	var rows []SKUWithRefs

	// join product_size_skus -> product_sizes (чтобы получить product_id)
	if err := r.DB.WithContext(ctx).
		Table("product_size_skus pssk").
		Select("pssk.sku as sku, ps.product_id as product_id, pssk.product_size_id as product_size_id").
		Joins("JOIN product_sizes ps ON ps.id = pssk.product_size_id").
		Where("ps.product_id IN ?", productIDs).
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_size_sku list by product ids: %w", err)
	}

	return rows, nil
}
