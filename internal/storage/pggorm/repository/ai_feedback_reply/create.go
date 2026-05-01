package ai_feedback_reply

import (
	"context"
	"fmt"
	"time"

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

// MarkWBSent проставляет wb_sent_at для записи — факт успешной отправки в WB.
func (r *Repository) MarkWBSent(ctx context.Context, replyID int, at time.Time) error {
	err := r.DB.WithContext(ctx).
		Model(&model.AIFeedbackReply{}).
		Where("id = ?", replyID).
		Update("wb_sent_at", at).Error
	if err != nil {
		return fmt.Errorf("ai_feedback_replies mark wb_sent: %w", err)
	}
	return nil
}

// ListPendingWBSend возвращает сгенерированные ответы, которые ещё не были отправлены в WB (wb_sent_at IS NULL).
// Делает JOIN с wb_feedbacks чтобы получить wb_feedback_id для вызова WB API.
func (r *Repository) ListPendingWBSend(ctx context.Context, limit int) ([]domain.AIFeedbackPendingSend, error) {
	if limit <= 0 {
		limit = 10
	}
	type row struct {
		ReplyID      int    `gorm:"column:reply_id"`
		ShopID       int    `gorm:"column:shop_id"`
		AnswerText   string `gorm:"column:answer_text"`
		WBFeedbackID string `gorm:"column:wb_feedback_id"`
	}
	var rows []row
	err := r.DB.WithContext(ctx).
		Table("ai_feedback_replies r").
		Select("r.id AS reply_id, r.shop_id, r.answer_text, f.wb_feedback_id").
		Joins("JOIN wb_feedbacks f ON f.id = r.feedback_id").
		Where("r.status = ? AND r.wb_sent_at IS NULL", domain.AIReplyStatusGenerated).
		Order("r.created_at ASC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("ai_feedback_replies list pending wb send: %w", err)
	}
	out := make([]domain.AIFeedbackPendingSend, len(rows))
	for i, row := range rows {
		out[i] = domain.AIFeedbackPendingSend{
			ReplyID:      row.ReplyID,
			ShopID:       row.ShopID,
			AnswerText:   row.AnswerText,
			WBFeedbackID: row.WBFeedbackID,
		}
	}
	return out, nil
}
