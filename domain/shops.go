package domain

import "time"

type Shop struct {
	ID          int
	UserID      int
	Marketplace string
	ExternalID  string
	Name        string
	APIToken    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
