package productsizesku

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) UpsertMany(ctx context.Context, productSizeID int, skus []string) error {
	if productSizeID == 0 || len(skus) == 0 {
		return nil
	}

	now := time.Now()
	rows := make([]model.ProductSizeSKU, 0, len(skus))
	for _, sku := range skus {
		if sku == "" {
			continue
		}
		rows = append(rows, model.ProductSizeSKU{
			ProductSizeID: productSizeID,
			SKU:           sku,
			CreatedAt:     now,
		})
	}
	if len(rows) == 0 {
		return nil
	}

	if err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "product_size_id"},
				{Name: "sku"},
			},
			DoNothing: true,
		}).
		Create(&rows).Error; err != nil {

		return fmt.Errorf("product_size_sku upsert many: %w", err)
	}

	return nil
}

func (r *Repository) ReplaceByProductSizeID(ctx context.Context, productSizeID int, skus []string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_size_id = ?", productSizeID).
			Delete(&model.ProductSizeSKU{}).Error; err != nil {
			return fmt.Errorf("product_size_sku replace delete: %w", err)
		}

		tmp := NewRepository(r.Log, tx)
		if err := tmp.UpsertMany(ctx, productSizeID, skus); err != nil {
			return err
		}

		return nil
	})
}
