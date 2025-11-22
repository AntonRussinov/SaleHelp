package product

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Create(ctx context.Context, p *domain.Product) error {
	m := toModel(p)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("product create: %w", err)
	}

	// пробрасываем обратно ID и timestamps в доменную модель
	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt

	return nil
}
