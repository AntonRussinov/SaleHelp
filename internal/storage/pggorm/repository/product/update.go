package product

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Update(ctx context.Context, p *domain.Product) error {
	m := toModel(p)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("product update: %w", err)
	}

	p.UpdatedAt = m.UpdatedAt

	return nil
}
