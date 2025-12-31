package domain

import "time"

type AlertTriggeredDelivery struct {
	ID          int
	TriggeredID int
	EndpointID  int

	Status    string // new/sent/failed
	SentAt    *time.Time
	Attempts  int
	LastError *string

	LeasedBy   *string
	LeaseUntil *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
