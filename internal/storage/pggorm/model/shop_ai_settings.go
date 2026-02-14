package model

import (
	"salehelp/domain"
	"time"
)

type ShopAISettings struct {
	ShopID      int       `gorm:"primaryKey;column:shop_id"`
	ToneOfVoice string    `gorm:"column:tone_of_voice;type:text;not null"`
	BrandName   string    `gorm:"column:brand_name;type:text"`
	CustomRules string    `gorm:"column:custom_rules;type:text"`
	IsActive    bool      `gorm:"column:is_active;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

func (m *ShopAISettings) ToDomain() *domain.ShopAISettings {
	if m == nil {
		return nil
	}
	return &domain.ShopAISettings{
		ShopID:      m.ShopID,
		ToneOfVoice: m.ToneOfVoice,
		BrandName:   m.BrandName,
		CustomRules: m.CustomRules,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (*ShopAISettings) TableName() string {
	return "shop_ai_settings"
}
