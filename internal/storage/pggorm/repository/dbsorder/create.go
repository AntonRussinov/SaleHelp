package dbsorder

import (
	"context"
	"fmt"

	"salehelp/domain"
)

// Create вставляет новый заказ DBS. При конфликте (shop_id, wb_order_id) ошибку не возвращаем —
// заказ уже есть, считаем это успехом (идемпотентность).
func (r *Repository) Create(ctx context.Context, o *domain.DBSOrder) error {
	m := toModel(o)
	err := r.DB.WithContext(ctx).Create(m).Error
	if err != nil {
		if isUniqueViolation(err) {
			return nil
		}
		return fmt.Errorf("dbs_order create: %w", err)
	}
	o.ID = m.ID
	return nil
}
