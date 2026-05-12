package subscription

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Create(ctx context.Context, s *domain.Subscription) error {
	m := toModel(s)
	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("subscription create: %w", err)
	}
	s.ID = m.ID
	s.CreatedAt = m.CreatedAt
	s.UpdatedAt = m.UpdatedAt
	return nil
}
