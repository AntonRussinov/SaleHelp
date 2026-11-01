package model

import (
	"salehelp/domain"
	"time"
)

type AITemplate struct {
	ID                int       `gorm:"primaryKey;column:id"`
	AISettingsID      int       `gorm:"column:ai_settings_id;not null"`
	Category          string    `gorm:"column:category;type:text;not null"`
	CustomPrompt      string    `gorm:"column:custom_prompt;type:text"`
	ToneVoiceOverride string    `gorm:"column:tone_voice_override;type:text"`
	Enabled           bool      `gorm:"column:enabled;not null"`
	CreatedAt         time.Time `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null"`
}

func (*AITemplate) TableName() string { return "ai_templates" }

func (m *AITemplate) ToDomain() *domain.AITemplate {
	if m == nil {
		return nil
	}
	return &domain.AITemplate{
		ID:                m.ID,
		AISettingsID:      m.AISettingsID,
		Category:          m.Category,
		CustomPrompt:      m.CustomPrompt,
		ToneVoiceOverride: m.ToneVoiceOverride,
		Enabled:           m.Enabled,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}
