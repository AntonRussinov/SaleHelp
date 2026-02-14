package ai

// FeedbackReplyInput — входные данные для генерации ответа на отзыв.
type FeedbackReplyInput struct {
	BrandName    string
	ToneInstruction string
	ProductName  string
	Rating       int
	Text         string
	Pros         string
	Cons         string
	CustomRules  string
}

// FeedbackReplyResult — результат генерации ответа (текст + метаданные).
type FeedbackReplyResult struct {
	Text       string
	TokensUsed int
	ModelUsed  string
}

// ModelConfig — имена моделей по рейтингу (rating 5 → mini, 1–4 → full).
type ModelConfig struct {
	ModelForPositive string // рейтинг 5
	ModelForOther    string // рейтинг 1–4
}

// DefaultModelConfig возвращает конфиг по умолчанию (реальные имена OpenAI).
func DefaultModelConfig() ModelConfig {
	return ModelConfig{
		ModelForPositive: "gpt-4o-mini",
		ModelForOther:    "gpt-4o",
	}
}
