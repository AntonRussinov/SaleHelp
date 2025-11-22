package airecommendation

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Create(ctx context.Context, d *domain.AIRecommendation) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("ai_recommendation create: %w", err)
	}

	d.ID = m.ID
	d.CreatedAt = m.CreatedAt

	return nil
}
