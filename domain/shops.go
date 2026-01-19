package domain

import "time"

type Shop struct {
	ID                  int
	UserID              int
	Marketplace         string
	ExternalID          string
	Name                string
	APIToken            string
	IsActive            bool
	LastSyncAt          *time.Time
	FirstSyncNotifiedAt *time.Time
	LastSyncStatusCode  int
	LastSyncError       string
	LastSyncFailedAt    *time.Time
	LastSyncOkAt        *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
