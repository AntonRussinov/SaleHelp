package wbfeedback

import (
	"context"
	"fmt"
	"time"

	"salehelp/internal/storage/pggorm/model"
)

// MarkNotified проставляет notified_at для указанных ID отзывов.
func (r *Repository) MarkNotified(ctx context.Context, ids []int, at time.Time) error {
	if len(ids) == 0 {
		return nil
	}
	err := r.DB.WithContext(ctx).
		Model(&model.WBFeedback{}).
		Where("id IN ?", ids).
		Update("notified_at", at).Error
	if err != nil {
		return fmt.Errorf("wb_feedbacks mark notified: %w", err)
	}
	return nil
}

// MarkReplied проставляет is_replied = true и replied_at для указанных ID отзывов.
func (r *Repository) MarkReplied(ctx context.Context, ids []int, at time.Time) error {
	if len(ids) == 0 {
		return nil
	}
	err := r.DB.WithContext(ctx).
		Model(&model.WBFeedback{}).
		Where("id IN ?", ids).
		Updates(map[string]any{"is_replied": true, "replied_at": at}).Error
	if err != nil {
		return fmt.Errorf("wb_feedbacks mark replied: %w", err)
	}
	return nil
}
