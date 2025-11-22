package domain

import "time"

type Product struct {
	ID          int
	ShopID      int
	ExternalID  string
	SKU         string
	Title       string
	URL         string
	ImageURL    string
	Brand       string
	Category    string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
