package domain

import "time"

type NotificationEndpoint struct {
	ID        int
	ShopID    int
	Channel   string // telegram/email/...
	Address   string
	Enabled   bool
	Label     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
