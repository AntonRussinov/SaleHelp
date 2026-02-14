package api

import "time"

// ToneOfVoice — допустимые значения тона ответов AI (для валидации и API).
const (
	ToneStandard  = "standard"
	ToneFriendly  = "friendly"
	ToneOfficial  = "official"
	ToneEmpathic  = "empathic"
	ToneMinimalist = "minimalist"
)

// ValidToneOfVoiceValues — список для валидации и документации.
var ValidToneOfVoiceValues = []string{ToneStandard, ToneFriendly, ToneOfficial, ToneEmpathic, ToneMinimalist}

// UpsertAISettingsRequest — запрос на создание/обновление настроек AI магазина.
type UpsertAISettingsRequest struct {
	// ToneOfVoice: standard | friendly | official | empathic | minimalist
	ToneOfVoice string `json:"toneOfVoice" validate:"required,oneof=standard friendly official empathic minimalist"`
	// BrandName — название бренда для промпта (опционально)
	BrandName string `json:"brandName" validate:"omitempty,max=500"`
	// CustomRules — доп. правила для промпта (опционально)
	CustomRules string `json:"customRules" validate:"omitempty,max=2000"`
	// IsActive — включена ли генерация AI-ответов для магазина
	IsActive bool `json:"isActive"`
}

// AISettingsResponse — ответ с настройками AI магазина.
type AISettingsResponse struct {
	ShopID      int       `json:"shopId"`
	ToneOfVoice string    `json:"toneOfVoice"`
	BrandName   string    `json:"brandName"`
	CustomRules string    `json:"customRules"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
