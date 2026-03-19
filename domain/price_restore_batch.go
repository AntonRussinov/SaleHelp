package domain

import "time"

type PriceRestoreBatch struct {
	ID         int
	CreatedAt  time.Time
	RestoredAt *time.Time
}
