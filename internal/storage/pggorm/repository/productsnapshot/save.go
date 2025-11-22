package productsnapshot

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Create(ctx context.Context, s *domain.ProductSnapshot) error {
	m := toModel(s)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("snapshot create: %w", err)
	}

	s.ID = m.ID
	s.CreatedAt = m.CreatedAt
	s.SnapshotTime = m.SnapshotTime

	return nil
}
