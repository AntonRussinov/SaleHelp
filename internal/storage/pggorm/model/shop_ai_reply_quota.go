package model

import (
	"salehelp/domain"
	"time"
)

type ShopAIReplyQuota struct {
	ID            int       `gorm:"primaryKey;column:id"`
	ShopID        int       `gorm:"column:shop_id;not null"`
	Month         time.Time `gorm:"column:month;type:date;not null"`
	IncludedLimit int       `gorm:"column:included_limit;not null"`
	IncludedUsed  int       `gorm:"column:included_used;not null"`
	PaidBalance   int       `gorm:"column:paid_balance;not null"`
	PaidUsed      int       `gorm:"column:paid_used;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
}

func (m *ShopAIReplyQuota) ToDomain() *domain.ShopAIReplyQuota {
	if m == nil {
		return nil
	}
	return &domain.ShopAIReplyQuota{
		ID:            m.ID,
		ShopID:        m.ShopID,
		Month:         m.Month,
		IncludedLimit: m.IncludedLimit,
		IncludedUsed:  m.IncludedUsed,
		PaidBalance:   m.PaidBalance,
		PaidUsed:      m.PaidUsed,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (*ShopAIReplyQuota) TableName() string {
	return "shop_ai_reply_quotas"
}
