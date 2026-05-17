package model

import (
	"salehelp/domain"
	"time"
)

type AIRule struct {
	ID           int       `gorm:"primaryKey;column:id"`
	AISettingsID int       `gorm:"column:ai_settings_id;not null"`
	RuleType     string    `gorm:"column:rule_type;type:text;not null"`
	RuleText     string    `gorm:"column:rule_text;type:text;not null"`
	OrderIndex   int       `gorm:"column:order_index;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
}

func (*AIRule) TableName() string { return "ai_rules" }

func (m *AIRule) ToDomain() *domain.AIRule {
	if m == nil {
		return nil
	}
	return &domain.AIRule{
		ID:           m.ID,
		AISettingsID: m.AISettingsID,
		RuleType:     m.RuleType,
		RuleText:     m.RuleText,
		OrderIndex:   m.OrderIndex,
		CreatedAt:    m.CreatedAt,
	}
}
