package ai_feedback_reply

import (
	"context"
	"fmt"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

// Create сохраняет запись о попытке/результате генерации ответа.
func (r *Repository) Create(ctx context.Context, d *domain.AIFeedbackReply) error {
	m := &model.AIFeedbackReply{
		FeedbackID:   d.FeedbackID,
		ShopID:       d.ShopID,
		AnswerText:   d.AnswerText,
		ModelUsed:    d.ModelUsed,
		Status:       d.Status,
		ErrorMessage: d.ErrorMessage,
		TokensUsed:   d.TokensUsed,
	}
	err := r.DB.WithContext(ctx).Create(m).Error
	if err != nil {
		return fmt.Errorf("ai_feedback_replies create: %w", err)
	}
	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	return nil
}
