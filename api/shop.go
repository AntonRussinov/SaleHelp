package api

import "time"

// CreateShopRequest тело POST /shops.
type CreateShopRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	APIToken    string `json:"apiToken" validate:"required"`
	Marketplace string `json:"marketplace" validate:"required,oneof=wildberries"`
}

// ShopResponse один магазин в ответе.
type ShopResponse struct {
	ID               int        `json:"id"`
	UserID           int        `json:"userId"`
	Marketplace      string     `json:"marketplace"`
	Name             string     `json:"name"`
	IsActive         bool       `json:"isActive"`
	LastSyncAt       *time.Time `json:"lastSyncAt,omitempty"`
	LastSyncStatusCode int      `json:"lastSyncStatusCode,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
}
