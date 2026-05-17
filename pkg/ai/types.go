package ai

// FeedbackReplyInput — входные данные для генерации ответа на отзыв.
type FeedbackReplyInput struct {
	BrandName            string
	BrandDescription     string // контекст о бренде (необязательно)
	ToneInstruction      string
	ProductName          string
	Rating               int
	Text                 string
	Pros                 string
	Cons                 string
	// OrderStatus — статус заказа из WB: "buyout" (выкуп), "returned" (товар возвращён).
	OrderStatus string
	// Addressing — обращение к покупателю: "vy" (вы) или "ty" (ты).
	Addressing string
	// EmojiMode — режим эмодзи: "none", "minimal", "expressive".
	EmojiMode string
	// ResponseLength — длина ответа: "short" (200), "medium" (400), "long" (600).
	ResponseLength string
	// Signature — подпись в конце ответа (необязательно).
	Signature string
	// CustomRules — список доп. правил бренда (каждый элемент — одно правило).
	CustomRules []string
	// ForbiddenWords — слова/фразы, которые нельзя упоминать в ответе.
	ForbiddenWords []string
	// CategoryCustomPrompt — доп. инструкции для конкретной категории (из AI-шаблонов).
	CategoryCustomPrompt string
}

// FeedbackReplyResult — результат генерации ответа (текст + метаданные).
type FeedbackReplyResult struct {
	Text         string
	InputTokens  int
	OutputTokens int
	TokensUsed   int // = InputTokens + OutputTokens
	ModelUsed    string
}

// ModelConfig — имена моделей по рейтингу отзыва.
type ModelConfig struct {
	ModelForPositive string // рейтинг 4–5 (позитив) — дешёвая
	ModelForMid      string // рейтинг 3 (нейтральный) — средняя
	ModelForOther    string // рейтинг 1–2 (негатив) — сложная
}

// SelectModel возвращает имя модели по рейтингу отзыва.
func (c ModelConfig) SelectModel(rating int) string {
	switch {
	case rating >= 4:
		return c.ModelForPositive
	case rating == 3:
		return c.ModelForMid
	default:
		return c.ModelForOther
	}
}

// DefaultModelConfig возвращает конфиг по умолчанию.
func DefaultModelConfig() ModelConfig {
	return ModelConfig{
		ModelForPositive: "gpt-4o-mini",
		ModelForMid:      "gpt-4o",
		ModelForOther:    "gpt-4o",
	}
}
