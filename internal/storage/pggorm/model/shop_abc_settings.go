package model

import (
	"salehelp/domain"
	"time"
)

type ShopABCSettings struct {
	ShopID             int       `gorm:"primaryKey;column:shop_id"`
	ABCWeeklyEnabled   bool      `gorm:"column:abc_weekly_enabled;not null;default:true"`
	ABCMonthlyEnabled  bool      `gorm:"column:abc_monthly_enabled;not null;default:true"`
	ABCAlertsEnabled   bool      `gorm:"column:abc_alerts_enabled;not null;default:true"`
	CreatedAt          time.Time `gorm:"column:created_at;not null"`
	UpdatedAt          time.Time `gorm:"column:updated_at;not null"`
}

func (m *ShopABCSettings) ToDomain() *domain.ShopABCSettings {
	if m == nil {
		return nil
	}
	return &domain.ShopABCSettings{
		ShopID:            m.ShopID,
		ABCWeeklyEnabled:  m.ABCWeeklyEnabled,
		ABCMonthlyEnabled: m.ABCMonthlyEnabled,
		ABCAlertsEnabled:  m.ABCAlertsEnabled,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func (*ShopABCSettings) TableName() string {
	return "shop_abc_settings"
}
