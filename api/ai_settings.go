package api

import "time"

// AI settings — допустимые значения полей.
const (
	AIResponseLengthShort  = "short"
	AIResponseLengthMedium = "medium"
	AIResponseLengthLong   = "long"

	AIAddressingVy = "vy"
	AIAddressingTy = "ty"

	AIEmojiModeNone       = "none"
	AIEmojiModeMinimal    = "minimal"
	AIEmojiModeExpressive = "expressive"

	AIRuleTypeRule         = "rule"
	AIRuleTypeForbiddenWord = "forbidden_word"

	AITemplateCategoryPositive = "positive"
	AITemplateCategoryNeutral  = "neutral"
	AITemplateCategoryNegative = "negative"
)

// UpsertAISettingsRequest — запрос на создание/обновление базовых настроек AI магазина.
type UpsertAISettingsRequest struct {
	IsActive           bool   `json:"isActive"           example:"true"`
	ToneOfVoice        string `json:"toneOfVoice"        validate:"required,oneof=standard friendly official empathic minimalist" example:"friendly"`
	ResponseLength     string `json:"responseLength"     validate:"required,oneof=short medium long"                              example:"medium"`
	Addressing         string `json:"addressing"         validate:"required,oneof=ty vy"                                          example:"vy"`
	EmojiMode          string `json:"emojiMode"          validate:"required,oneof=none minimal expressive"                        example:"minimal"`
	BrandName          string `json:"brandName"          validate:"omitempty,max=100"                                            example:"МойБренд"`
	BrandDescription   string `json:"brandDescription"   validate:"omitempty,max=500"                                            example:"Производим эко-товары для дома"`
	RepresentativeName string `json:"representativeName" validate:"omitempty,max=50"                                             example:"Иван"`
	Signature          string `json:"signature"          validate:"omitempty,max=200"                                            example:"С уважением, команда МойБренд"`
}

// AISettingsResponse — полные настройки AI магазина с правилами и шаблонами.
type AISettingsResponse struct {
	ID                 int                  `json:"id"`
	ShopID             int                  `json:"shopId"`
	IsActive           bool                 `json:"isActive"`
	ToneOfVoice        string               `json:"toneOfVoice"`
	ResponseLength     string               `json:"responseLength"`
	Addressing         string               `json:"addressing"`
	EmojiMode          string               `json:"emojiMode"`
	BrandName          string               `json:"brandName"`
	BrandDescription   string               `json:"brandDescription"`
	RepresentativeName string               `json:"representativeName"`
	Signature          string               `json:"signature"`
	Rules              []AIRuleResponse     `json:"rules"`
	Templates          []AITemplateResponse `json:"templates"`
	CreatedAt          time.Time            `json:"createdAt"`
	UpdatedAt          time.Time            `json:"updatedAt"`
}

// AIRuleResponse — одно правило бренда или запрещённое слово.
type AIRuleResponse struct {
	ID         int    `json:"id"`
	RuleType   string `json:"ruleType"   example:"rule"`
	RuleText   string `json:"ruleText"   example:"Всегда упоминай гарантию 1 год"`
	OrderIndex int    `json:"orderIndex"`
}

// CreateAIRuleRequest — запрос на добавление правила.
type CreateAIRuleRequest struct {
	RuleType   string `json:"ruleType"   validate:"required,oneof=rule forbidden_word" example:"rule"`
	RuleText   string `json:"ruleText"   validate:"required,max=200"                   example:"Всегда упоминай гарантию 1 год"`
	OrderIndex int    `json:"orderIndex" validate:"min=0"                              example:"0"`
}

// UpdateAIRuleRequest — запрос на обновление правила.
type UpdateAIRuleRequest struct {
	RuleType   string `json:"ruleType"   validate:"required,oneof=rule forbidden_word"`
	RuleText   string `json:"ruleText"   validate:"required,max=200"`
	OrderIndex int    `json:"orderIndex" validate:"min=0"`
}

// AITemplateResponse — настройки для конкретной категории отзыва.
type AITemplateResponse struct {
	ID                int    `json:"id"`
	Category          string `json:"category"          example:"positive"`
	CustomPrompt      string `json:"customPrompt"      example:"Упомяни программу лояльности"`
	ToneVoiceOverride string `json:"toneVoiceOverride" example:""`
	Enabled           bool   `json:"enabled"`
}

// UpsertAITemplateRequest — запрос на обновление шаблона категории.
type UpsertAITemplateRequest struct {
	CustomPrompt      string `json:"customPrompt"      validate:"omitempty,max=500"`
	ToneVoiceOverride string `json:"toneVoiceOverride" validate:"omitempty,oneof=standard friendly official empathic minimalist"`
	Enabled           *bool  `json:"enabled"`
}
