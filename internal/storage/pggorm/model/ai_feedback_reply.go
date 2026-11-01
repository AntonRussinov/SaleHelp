package model

import (
	"salehelp/domain"
	"time"
)

type AIFeedbackReply struct {
	ID                  int        `gorm:"primaryKey;column:id"`
	FeedbackID          int        `gorm:"column:feedback_id;not null"`
	ShopID              int        `gorm:"column:shop_id;not null"`
	AnswerText          string     `gorm:"column:answer_text;type:text"`
	ModelUsed           string     `gorm:"column:model_used;type:text;not null"`
	Status              string     `gorm:"column:status;type:text;not null"`
	ErrorMessage        string     `gorm:"column:error_message;type:text"`
	InputTokens         int        `gorm:"column:input_tokens;not null"`
	OutputTokens        int        `gorm:"column:output_tokens;not null"`
	TokensUsed          int        `gorm:"column:tokens_used;not null"`
	AITier              string     `gorm:"column:ai_tier;type:text;not null"`
	EstimatedCostKopeks int        `gorm:"column:estimated_cost_kopeks;not null"`
	WBSentAt            *time.Time `gorm:"column:wb_sent_at"`
	CreatedAt           time.Time  `gorm:"column:created_at;not null"`
}

func (m *AIFeedbackReply) ToDomain() *domain.AIFeedbackReply {
	if m == nil {
		return nil
	}
	return &domain.AIFeedbackReply{
		ID:                  m.ID,
		FeedbackID:          m.FeedbackID,
		ShopID:              m.ShopID,
		AnswerText:          m.AnswerText,
		ModelUsed:           m.ModelUsed,
		Status:              m.Status,
		ErrorMessage:        m.ErrorMessage,
		InputTokens:         m.InputTokens,
		OutputTokens:        m.OutputTokens,
		TokensUsed:          m.TokensUsed,
		AITier:              m.AITier,
		EstimatedCostKopeks: m.EstimatedCostKopeks,
		WBSentAt:            m.WBSentAt,
		CreatedAt:           m.CreatedAt,
	}
}

func (*AIFeedbackReply) TableName() string {
	return "ai_feedback_replies"
}
