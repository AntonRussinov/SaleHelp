package shop

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) Update(ctx context.Context, s *domain.Shop) error {
	m := toModel(s)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("shop update: %w", err)
	}

	s.UpdatedAt = m.UpdatedAt

	return nil
}

func (r *Repository) SetActive(
	ctx context.Context,
	shopID int,
	active bool,
) error {
	err := r.DB.WithContext(ctx).
		Model(&model.Shop{}).
		Where("id = ?", shopID).
		Update("is_active", active).Error

	if err != nil {
		return fmt.Errorf("shop set active=%v: %w", active, err)
	}
	return nil
}
