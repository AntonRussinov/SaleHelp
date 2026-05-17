package domain

import "time"

const (
	AIRuleTypeRule         = "rule"
	AIRuleTypeForbiddenWord = "forbidden_word"

	AITemplateCategoryPositive = "positive"
	AITemplateCategoryNeutral  = "neutral"
	AITemplateCategoryNegative = "negative"

	AIResponseLengthShort  = "short"
	AIResponseLengthMedium = "medium"
	AIResponseLengthLong   = "long"

	AIAddressingVy = "vy"
	AIAddressingTy = "ty"

	AIEmojiModeNone       = "none"
	AIEmojiModeMinimal    = "minimal"
	AIEmojiModeExpressive = "expressive"
)

// AISettings — настройки AI-ответов для магазина (новая схема, заменяет shop_ai_settings).
type AISettings struct {
	ID                 int
	ShopID             int
	UserID             int
	IsActive           bool
	ToneOfVoice        string // standard | friendly | official | empathic | minimalist
	ResponseLength     string // short | medium | long
	Addressing         string // ty | vy
	EmojiMode          string // none | minimal | expressive
	BrandName          string
	BrandDescription   string
	RepresentativeName string
	Signature          string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// AIRule — одно правило бренда или запрещённое слово.
type AIRule struct {
	ID           int
	AISettingsID int
	RuleType     string // rule | forbidden_word
	RuleText     string
	OrderIndex   int
	CreatedAt    time.Time
}

// AITemplate — кастомные инструкции и настройки для конкретной категории отзыва.
type AITemplate struct {
	ID                int
	AISettingsID      int
	Category          string // positive | neutral | negative
	CustomPrompt      string
	ToneVoiceOverride string // если задан — перекрывает базовый ToneOfVoice
	Enabled           bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// AIExample — few-shot пример ответа (Pro, пока не используется в промпте).
type AIExample struct {
	ID           int
	AISettingsID int
	ReviewText   string
	IdealResponse string
	Category     string // positive | neutral | negative | ""
	CreatedAt    time.Time
}

// AISettingsBundle — всё нужное для генерации промпта одним запросом.
type AISettingsBundle struct {
	Settings *AISettings
	Rules    []*AIRule
	Template *AITemplate // шаблон для нужной категории (может быть nil)
}

// CategoryFromRating возвращает категорию шаблона по рейтингу отзыва.
func CategoryFromRating(rating int) string {
	switch {
	case rating >= 4:
		return AITemplateCategoryPositive
	case rating == 3:
		return AITemplateCategoryNeutral
	default:
		return AITemplateCategoryNegative
	}
}
