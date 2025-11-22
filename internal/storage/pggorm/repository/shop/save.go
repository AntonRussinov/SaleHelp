package shop

import (
	"context"
	"fmt"
	"salehelp/domain"
)

// Create создаёт магазин
func (r *Repository) Create(ctx context.Context, s *domain.Shop) error {
	m := toModel(s)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("shop create: %w", err)
	}

	s.ID = m.ID
	s.CreatedAt = m.CreatedAt
	s.UpdatedAt = m.UpdatedAt

	return nil
}
