package model

import (
	"salehelp/domain"
	"time"
)

type ShopNotificationEndpoint struct {
	ID      int     `gorm:"primaryKey;column:id"`
	ShopID  int     `gorm:"column:shop_id;not null"`
	Channel string  `gorm:"column:channel;type:text;not null"`
	Address string  `gorm:"column:address;type:text;not null"`
	Enabled bool    `gorm:"column:enabled;not null;default:true"`
	Label   *string `gorm:"column:label;type:text"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func ShopNotificationEndpointColumns() []string {
	return []string{
		"id",
		"shop_id",
		"channel",
		"address",
		"enabled",
		"label",
		"created_at",
		"updated_at",
	}
}

func (m *ShopNotificationEndpoint) ToDomain() *domain.NotificationEndpoint {
	if m == nil {
		return nil
	}

	return &domain.NotificationEndpoint{
		ID:        m.ID,
		ShopID:    m.ShopID,
		Channel:   m.Channel,
		Address:   m.Address,
		Enabled:   m.Enabled,
		Label:     m.Label,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (*ShopNotificationEndpoint) TableName() string {
	return "shop_notification_endpoints"
}
