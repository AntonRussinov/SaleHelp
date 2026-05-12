package model

import (
	"salehelp/domain"
	"time"
)

type Subscription struct {
	ID           int        `gorm:"primaryKey;column:id"`
	UserID       int        `gorm:"column:user_id;not null"`
	Tier         string     `gorm:"column:tier;type:text;not null"`
	Status       string     `gorm:"column:status;type:text;not null"`
	PeriodMonths int        `gorm:"column:period_months;not null"`
	StartedAt    time.Time  `gorm:"column:started_at;not null"`
	ExpiresAt    time.Time  `gorm:"column:expires_at;not null"`
	TrialEndsAt  *time.Time `gorm:"column:trial_ends_at"`
	AIQuota      int        `gorm:"column:ai_quota;not null"`
	AIUsed       int        `gorm:"column:ai_used;not null"`
	AIResetsAt   time.Time  `gorm:"column:ai_resets_at;not null"`
	MaxShops     int        `gorm:"column:max_shops;not null"`
	ExtraShops   int        `gorm:"column:extra_shops;not null"`
	PriceMinor   int        `gorm:"column:price_minor;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (*Subscription) TableName() string { return "subscriptions" }

func (s *Subscription) ToDomain() *domain.Subscription {
	if s == nil {
		return nil
	}
	return &domain.Subscription{
		ID:           s.ID,
		UserID:       s.UserID,
		Tier:         s.Tier,
		Status:       s.Status,
		PeriodMonths: s.PeriodMonths,
		StartedAt:    s.StartedAt,
		ExpiresAt:    s.ExpiresAt,
		TrialEndsAt:  s.TrialEndsAt,
		AIQuota:      s.AIQuota,
		AIUsed:       s.AIUsed,
		AIResetsAt:   s.AIResetsAt,
		MaxShops:     s.MaxShops,
		ExtraShops:   s.ExtraShops,
		PriceMinor:   s.PriceMinor,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}
