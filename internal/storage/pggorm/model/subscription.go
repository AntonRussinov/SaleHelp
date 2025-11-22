package model

import (
	"salehelp/domain"
	"time"
)

type Subscription struct {
	ID                     int       `gorm:"primaryKey;column:id"`
	UserID                 int       `gorm:"column:user_id;not null"`
	Plan                   string    `gorm:"column:plan;type:text;not null"`
	Status                 string    `gorm:"column:status;type:text;not null"`
	ValidFrom              time.Time `gorm:"column:valid_from;not null"`
	ValidTo                time.Time `gorm:"column:valid_to;not null"`
	ExternalCustomerID     string    `gorm:"column:external_customer_id;type:text"`
	ExternalSubscriptionID string    `gorm:"column:external_subscription_id;type:text"`
	CreatedAt              time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt              time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func SubscriptionColumns() []string {
	return []string{
		"id",
		"user_id",
		"plan",
		"status",
		"valid_from",
		"valid_to",
		"external_customer_id",
		"external_subscription_id",
		"created_at",
		"updated_at",
	}
}

func (s *Subscription) ToDomain() *domain.Subscription {
	if s == nil {
		return nil
	}

	return &domain.Subscription{
		ID:                     s.ID,
		UserID:                 s.UserID,
		Plan:                   s.Plan,
		Status:                 s.Status,
		ValidFrom:              s.ValidFrom,
		ValidTo:                s.ValidTo,
		ExternalCustomerID:     s.ExternalCustomerID,
		ExternalSubscriptionID: s.ExternalSubscriptionID,
		CreatedAt:              s.CreatedAt,
		UpdatedAt:              s.UpdatedAt,
	}
}

func (*Subscription) TableName() string {
	return "subscriptions"
}
