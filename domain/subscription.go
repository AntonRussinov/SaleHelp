package domain

import "time"

type Subscription struct {
	ID              int
	ShopID          int
	Plan            string
	Status          string
	ValidFrom       time.Time
	ValidTo         time.Time
	PeriodMonths    int
	BaseAmountMinor int64
	CurrencyCode    string
	PaidAt          *time.Time
	Source          string
	Comment         *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
