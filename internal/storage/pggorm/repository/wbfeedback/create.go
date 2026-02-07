package wbfeedback

import (
	"context"
	"fmt"

	"salehelp/domain"

	"gorm.io/gorm/clause"
)

// Create вставляет новый отзыв. При конфликте (shop_id, wb_feedback_id) ничего не делаем.
func (r *Repository) Create(ctx context.Context, f *domain.WBFeedback) error {
	m := toModel(f)
	err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "shop_id"}, {Name: "wb_feedback_id"}},
			DoNothing: true,
		}).
		Create(m).Error
	if err != nil {
		return fmt.Errorf("wb_feedback create: %w", err)
	}
	if m.ID != 0 {
		f.ID = m.ID
	}
	return nil
}
