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
	ToneOfVoice string `json:"toneOfVoice" validate:"required,oneof=standard friendly official empathic minimalist" enums:"standard,friendly,official,empathic,minimalist" example:"friendly"`
	// BrandName — название бренда для промпта (опционально)
	BrandName string `json:"brandName" validate:"omitempty,max=500" example:"МойБренд"`
	// CustomRules — доп. правила для промпта (опционально)
	CustomRules string `json:"customRules" validate:"omitempty,max=2000" example:"Всегда упоминай гарантию. Не обещай скидки."`
	// IsActive — включена ли генерация AI-ответов для магазина
	IsActive bool `json:"isActive" example:"true"`
}

// AISettingsResponse — ответ с настройками AI магазина.
type AISettingsResponse struct {
	ShopID      int       `json:"shopId"      example:"1"`
	ToneOfVoice string    `json:"toneOfVoice" example:"friendly"`
	BrandName   string    `json:"brandName"   example:"МойБренд"`
	CustomRules string    `json:"customRules" example:"Всегда упоминай гарантию."`
	IsActive    bool      `json:"isActive"    example:"true"`
	CreatedAt   time.Time `json:"createdAt"   example:"2024-01-15T10:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt"   example:"2024-01-15T10:00:00Z"`
}
