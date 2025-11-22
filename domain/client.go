package domain

import "time"

type Client struct {
	ID        int
	Name      string
	WBAPIKey  string
	CreatedAt time.Time
}
