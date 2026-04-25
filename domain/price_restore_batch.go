package domain

import "time"

const (
	RestoreStatusPending   = "pending"
	RestoreStatusRestoring = "restoring"
	RestoreStatusRestored  = "restored"
	RestoreStatusFailed    = "failed"
)

type PriceRestoreBatch struct {
	ID               int
	ShopID           int
	RestoreStatus    string
	RestoreStartedAt *time.Time
	CreatedAt        time.Time
	RestoredAt       *time.Time
}
