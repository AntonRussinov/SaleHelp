package wbquestion

import (
	"context"
	"fmt"

	"salehelp/domain"

	"gorm.io/gorm/clause"
)

// Create вставляет новый вопрос. При конфликте (shop_id, wb_question_id) ничего не делаем.
func (r *Repository) Create(ctx context.Context, q *domain.WBQuestion) error {
	m := toModel(q)
	err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "shop_id"}, {Name: "wb_question_id"}},
			DoNothing: true,
		}).
		Create(m).Error
	if err != nil {
		return fmt.Errorf("wb_question create: %w", err)
	}
	if m.ID != 0 {
		q.ID = m.ID
	}
	return nil
}
