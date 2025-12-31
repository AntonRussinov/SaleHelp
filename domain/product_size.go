package domain

import "time"

type ProductSize struct {
	ID           int
	ProductID    int
	ExternalID   string  // chrtID
	TechSizeName *string // nullable (может не прийти)

	CreatedAt time.Time
	UpdatedAt time.Time
}
