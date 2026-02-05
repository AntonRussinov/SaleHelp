package fbsorder

import (
	"context"
	"fmt"

	"salehelp/domain"

	"gorm.io/gorm/clause"
)

// Create вставляет новый заказ. При конфликте (shop_id, wb_order_id) ничего не делаем —
// ON CONFLICT DO NOTHING, чтобы не получать ошибку от БД и не засорять логи.
func (r *Repository) Create(ctx context.Context, o *domain.FBSOrder) error {
	m := toModel(o)
	err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "shop_id"}, {Name: "wb_order_id"}},
			DoNothing: true,
		}).
		Create(m).Error
	if err != nil {
		return fmt.Errorf("fbs_order create: %w", err)
	}
	if m.ID != 0 {
		o.ID = m.ID
	}
	return nil
}
