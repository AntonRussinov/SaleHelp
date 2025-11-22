package alert

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Create(ctx context.Context, a *domain.Alert) error {
	m := toModel(a)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("alert create: %w", err)
	}

	a.ID = m.ID
	a.TriggeredAt = m.TriggeredAt

	return nil
}
