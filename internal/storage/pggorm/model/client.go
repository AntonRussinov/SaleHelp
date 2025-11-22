package model

import (
	"salehelp/domain"
	"time"
)

type Client struct {
	ID        int       `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	WBAPIKey  string    `gorm:"column:wb_api_key;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func ClientColumns() []string {
	return []string{
		"id",
		"name",
		"wb_api_key",
		"created_at",
	}
}

func (c *Client) ToDomain() *domain.Client {
	if c == nil {
		return nil
	}

	return &domain.Client{
		ID:        c.ID,
		Name:      c.Name,
		WBAPIKey:  c.WBAPIKey,
		CreatedAt: c.CreatedAt,
	}
}

func (*Client) TableName() string {
	return "clients"
}
