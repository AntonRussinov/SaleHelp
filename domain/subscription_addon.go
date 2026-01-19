package domain

import "time"

type SubscriptionAddon struct {
	ID             int
	SubscriptionID int
	Code           string
	Quantity       int
	AmountMinor    int64
	CurrencyCode   string
	Meta           []byte // json
	CreatedAt      time.Time
}
