package domain

import "time"

type AlertTriggered struct {
	ID         int
	AlertID    int
	ProductID  int
	SnapshotID int
	ShopID     int
	Message    string
	CreatedAt  time.Time
}
