package productsizepricesnapshot

import (
	"context"
	"fmt"
	"salehelp/domain"

	"gorm.io/gorm/clause"
)

func (r *Repository) Upsert(ctx context.Context, d *domain.ProductSizePriceSnapshot) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "product_snapshot_id"},
				{Name: "product_size_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"base_price",
				"final_price",
				"club_price",
				"discount",
				"club_discount",
				"extra",
			}),
		}).
		Create(m).Error; err != nil {

		return fmt.Errorf("product_size_price_snapshot upsert: %w", err)
	}

	d.ID = m.ID
	return nil
}
