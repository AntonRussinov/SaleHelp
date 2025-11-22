package shop

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Update(ctx context.Context, s *domain.Shop) error {
	m := toModel(s)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("shop update: %w", err)
	}

	s.UpdatedAt = m.UpdatedAt

	return nil
}
