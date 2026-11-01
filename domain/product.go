package domain

import "time"

type Product struct {
	ID          int       `json:"id"`
	ShopID      int       `json:"shopId"`
	ExternalID  string    `json:"externalId"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"imageUrl"`
	Brand       string    `json:"brand"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
