package model

import "time"

type ShopNotificationSetting struct {
	ID        int       `json:"id" gorm:"primaryKey;column:id"`
	ShopID    int       `json:"shopId" gorm:"column:shop_id;not null;index:idx_shop_notification_settings_shop_event,unique"`
	EventKey  string    `json:"eventKey" gorm:"column:event_key;type:text;not null;index:idx_shop_notification_settings_shop_event,unique"`
	Enabled   bool      `json:"enabled" gorm:"column:enabled;not null;default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (*ShopNotificationSetting) TableName() string {
	return "shop_notification_settings"
}

func (s *ShopNotificationSetting) GetEventKey() string {
	if s == nil {
		return ""
	}
	return s.EventKey
}

func (s *ShopNotificationSetting) GetEnabled() bool {
	if s == nil {
		return false
	}
	return s.Enabled
}
