package user

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Update(ctx context.Context, u *domain.User) error {
	m := toModel(u)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("user update: %w", err)
	}

	u.UpdatedAt = m.UpdatedAt

	return nil
}
