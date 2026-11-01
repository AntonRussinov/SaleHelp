package wbquestion

import (
	"context"
	"errors"
	"fmt"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

// GetByWBQuestionID returns a question by its WB API ID within a shop, or nil if not found.
func (r *Repository) GetByWBQuestionID(ctx context.Context, shopID int, wbQuestionID string) (*domain.WBQuestion, error) {
	var m model.WBQuestion
	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND wb_question_id = ?", shopID, wbQuestionID).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("wb_question get by wb_question_id: %w", err)
	}
	return m.ToDomain(), nil
}
