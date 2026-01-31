package saleevent

import (
	"context"
	"fmt"

	"salehelp/domain"
)

// Create вставляет новую запись продажи/возврата. При конфликте (shop_id, srid) ошибку не возвращаем —
// запись уже есть, считаем это успехом (идемпотентность).
func (r *Repository) Create(ctx context.Context, e *domain.SaleEvent) error {
	m := toModel(e)
	err := r.DB.WithContext(ctx).Create(m).Error
	if err != nil {
		if isUniqueViolation(err) {
			return nil
		}
		return fmt.Errorf("sale_event create: %w", err)
	}
	e.ID = m.ID
	return nil
}
