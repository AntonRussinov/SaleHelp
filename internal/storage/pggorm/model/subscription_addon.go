package model

import (
	"salehelp/domain"
	"time"
)

type SubscriptionAddon struct {
	ID             int       `gorm:"primaryKey;column:id"`
	SubscriptionID int       `gorm:"column:subscription_id;not null"`
	Code           string    `gorm:"column:code;type:text;not null"`
	Quantity       int       `gorm:"column:quantity;not null"`
	AmountMinor    int64     `gorm:"column:amount_minor;not null"`
	CurrencyCode   string    `gorm:"column:currency_code;type:char(3);not null"`
	Meta           []byte    `gorm:"column:meta;type:jsonb;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (*SubscriptionAddon) TableName() string { return "subscription_addons" }

func (a *SubscriptionAddon) ToDomain() *domain.SubscriptionAddon {
	if a == nil {
		return nil
	}
	return &domain.SubscriptionAddon{
		ID:             a.ID,
		SubscriptionID: a.SubscriptionID,
		Code:           a.Code,
		Quantity:       a.Quantity,
		AmountMinor:    a.AmountMinor,
		CurrencyCode:   a.CurrencyCode,
		Meta:           a.Meta,
		CreatedAt:      a.CreatedAt,
	}
}
