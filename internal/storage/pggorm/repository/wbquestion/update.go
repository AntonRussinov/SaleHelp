package wbquestion

import (
	"context"
	"fmt"
	"time"

	"salehelp/internal/storage/pggorm/model"
)

// MarkNotified проставляет notified_at для указанных ID вопросов.
func (r *Repository) MarkNotified(ctx context.Context, ids []int, at time.Time) error {
	if len(ids) == 0 {
		return nil
	}
	err := r.DB.WithContext(ctx).
		Model(&model.WBQuestion{}).
		Where("id IN ?", ids).
		Update("notified_at", at).Error
	if err != nil {
		return fmt.Errorf("wb_questions mark notified: %w", err)
	}
	return nil
}
