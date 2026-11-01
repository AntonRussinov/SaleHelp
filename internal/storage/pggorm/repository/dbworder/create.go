package dbworder

import (
	"context"
	"fmt"

	"salehelp/domain"

	"gorm.io/gorm/clause"
)

// Create вставляет новый заказ DBW. При конфликте (shop_id, wb_order_id) — ON CONFLICT DO NOTHING.
func (r *Repository) Create(ctx context.Context, o *domain.DBWOrder) error {
	m := toModel(o)
	err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "shop_id"}, {Name: "wb_order_id"}},
			DoNothing: true,
		}).
		Create(m).Error
	if err != nil {
		return fmt.Errorf("dbw_order create: %w", err)
	}
	if m.ID != 0 {
		o.ID = m.ID
	}
	return nil
}
