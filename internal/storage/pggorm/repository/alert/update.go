package alert

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Update(ctx context.Context, a *domain.Alert) error {
	m := toModel(a)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("alert update: %w", err)
	}

	// тут особо нечего синхронизировать обратно, разве что ResolvedAt, если меняли
	a.ResolvedAt = m.ResolvedAt

	return nil
}
