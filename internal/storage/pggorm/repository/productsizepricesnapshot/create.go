package productsizepricesnapshot

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Create(ctx context.Context, d *domain.ProductSizePriceSnapshot) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("product_size_price_snapshot create: %w", err)
	}

	d.ID = m.ID
	return nil
}
