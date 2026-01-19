package domain

import "time"

type SubscriptionNotificationEvent struct {
	ID             int
	SubscriptionID int
	Event          string // expiring_3d | expiring_0d | expired_1d
	SentAt         time.Time
}
