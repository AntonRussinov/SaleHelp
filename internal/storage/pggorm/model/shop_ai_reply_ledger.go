package model

import (
	"salehelp/domain"
	"time"
)

type ShopAIReplyLedger struct {
	ID         int        `gorm:"primaryKey;column:id"`
	ShopID     int        `gorm:"column:shop_id;not null"`
	FeedbackID *int       `gorm:"column:feedback_id"`
	Month      time.Time  `gorm:"column:month;type:date;not null"`
	Source     string     `gorm:"column:source;type:text;not null"`
	Delta      int        `gorm:"column:delta;not null"`
	Reason     string     `gorm:"column:reason;type:text;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null"`
}

func (m *ShopAIReplyLedger) ToDomain() *domain.ShopAIReplyLedger {
	if m == nil {
		return nil
	}
	return &domain.ShopAIReplyLedger{
		ID:         m.ID,
		ShopID:     m.ShopID,
		FeedbackID: m.FeedbackID,
		Month:      m.Month,
		Source:     m.Source,
		Delta:      m.Delta,
		Reason:     m.Reason,
		CreatedAt:  m.CreatedAt,
	}
}

func (*ShopAIReplyLedger) TableName() string {
	return "shop_ai_reply_ledger"
}
