package domain

import "time"

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusExpired  SubscriptionStatus = "expired"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
)

type Subscription struct {
	ID              int64
	ShopID          int64
	Plan            string
	Status          SubscriptionStatus
	ValidFrom       time.Time
	ValidTo         time.Time
	PeriodMonths    int16
	BaseAmountMinor int64
	CurrencyCode    string
	PaidAt          *time.Time
	Source          string
	Comment         *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
