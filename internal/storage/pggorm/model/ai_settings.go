package model

import (
	"salehelp/domain"
	"time"
)

type AISettings struct {
	ID                 int        `gorm:"primaryKey;column:id"`
	ShopID             int        `gorm:"column:shop_id;not null;uniqueIndex"`
	UserID             int        `gorm:"column:user_id;not null"`
	IsActive           bool       `gorm:"column:is_active;not null"`
	ToneOfVoice        string     `gorm:"column:tone_of_voice;type:text;not null"`
	ResponseLength     string     `gorm:"column:response_length;type:text;not null"`
	Addressing         string     `gorm:"column:addressing;type:text;not null"`
	EmojiMode          string     `gorm:"column:emoji_mode;type:text;not null"`
	BrandName          string     `gorm:"column:brand_name;type:text"`
	BrandDescription   string     `gorm:"column:brand_description;type:text"`
	RepresentativeName string     `gorm:"column:representative_name;type:text"`
	Signature          string     `gorm:"column:signature;type:text"`
	CreatedAt          time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;not null"`
}

func (*AISettings) TableName() string { return "ai_settings" }

func (m *AISettings) ToDomain() *domain.AISettings {
	if m == nil {
		return nil
	}
	return &domain.AISettings{
		ID:                 m.ID,
		ShopID:             m.ShopID,
		UserID:             m.UserID,
		IsActive:           m.IsActive,
		ToneOfVoice:        m.ToneOfVoice,
		ResponseLength:     m.ResponseLength,
		Addressing:         m.Addressing,
		EmojiMode:          m.EmojiMode,
		BrandName:          m.BrandName,
		BrandDescription:   m.BrandDescription,
		RepresentativeName: m.RepresentativeName,
		Signature:          m.Signature,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func AISettingsFromDomain(d *domain.AISettings) *AISettings {
	return &AISettings{
		ID:                 d.ID,
		ShopID:             d.ShopID,
		UserID:             d.UserID,
		IsActive:           d.IsActive,
		ToneOfVoice:        d.ToneOfVoice,
		ResponseLength:     d.ResponseLength,
		Addressing:         d.Addressing,
		EmojiMode:          d.EmojiMode,
		BrandName:          d.BrandName,
		BrandDescription:   d.BrandDescription,
		RepresentativeName: d.RepresentativeName,
		Signature:          d.Signature,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}
}
