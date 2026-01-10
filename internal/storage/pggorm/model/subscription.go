package model

import (
	"salehelp/domain"
	"time"
)

type Subscription struct {
	ID              int        `gorm:"primaryKey;column:id"`
	ShopID          int        `gorm:"column:shop_id;not null"`
	Plan            string     `gorm:"column:plan;type:text;not null"`
	Status          string     `gorm:"column:status;type:text;not null"`
	ValidFrom       time.Time  `gorm:"column:valid_from;not null"`
	ValidTo         time.Time  `gorm:"column:valid_to;not null"`
	PeriodMonths    int        `gorm:"column:period_months;not null"`
	BaseAmountMinor int64      `gorm:"column:base_amount_minor;not null"`
	CurrencyCode    string     `gorm:"column:currency_code;type:char(3);not null"`
	PaidAt          *time.Time `gorm:"column:paid_at"`
	Source          string     `gorm:"column:source;type:text;not null"`
	Comment         *string    `gorm:"column:comment;type:text"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (*Subscription) TableName() string { return "subscriptions" }

func (s *Subscription) ToDomain() *domain.Subscription {
	if s == nil {
		return nil
	}
	return &domain.Subscription{
		ID:              s.ID,
		ShopID:          s.ShopID,
		Plan:            s.Plan,
		Status:          s.Status,
		ValidFrom:       s.ValidFrom,
		ValidTo:         s.ValidTo,
		PeriodMonths:    s.PeriodMonths,
		BaseAmountMinor: s.BaseAmountMinor,
		CurrencyCode:    s.CurrencyCode,
		PaidAt:          s.PaidAt,
		Source:          s.Source,
		Comment:         s.Comment,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}
