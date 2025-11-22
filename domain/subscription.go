package domain

import "time"

type Subscription struct {
	ID                     int
	UserID                 int
	Plan                   string
	Status                 string
	ValidFrom              time.Time
	ValidTo                time.Time
	ExternalCustomerID     string
	ExternalSubscriptionID string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
