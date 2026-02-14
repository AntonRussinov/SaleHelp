package model

import (
	"salehelp/domain"
	"time"
)

type ShopAIReplyPurchase struct {
	ID           int        `gorm:"primaryKey;column:id"`
	ShopID       int        `gorm:"column:shop_id;not null"`
	PackageName  string     `gorm:"column:package_name;type:text;not null"`
	AnswersCount int        `gorm:"column:answers_count;not null"`
	PriceRub     int        `gorm:"column:price_rub;not null"`
	PaymentRef   string     `gorm:"column:payment_ref;type:text"`
	Status       string     `gorm:"column:status;type:text;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null"`
	PaidAt       *time.Time `gorm:"column:paid_at"`
}

func (m *ShopAIReplyPurchase) ToDomain() *domain.ShopAIReplyPurchase {
	if m == nil {
		return nil
	}
	return &domain.ShopAIReplyPurchase{
		ID:           m.ID,
		ShopID:       m.ShopID,
		PackageName:  m.PackageName,
		AnswersCount: m.AnswersCount,
		PriceRub:     m.PriceRub,
		PaymentRef:   m.PaymentRef,
		Status:       m.Status,
		CreatedAt:    m.CreatedAt,
		PaidAt:       m.PaidAt,
	}
}

func (*ShopAIReplyPurchase) TableName() string {
	return "shop_ai_reply_purchases"
}
