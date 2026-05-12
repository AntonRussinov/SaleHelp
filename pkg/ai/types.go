package ai

// FeedbackReplyInput — входные данные для генерации ответа на отзыв.
type FeedbackReplyInput struct {
	BrandName       string
	ToneInstruction string
	ProductName     string
	Rating          int
	Text            string
	Pros            string
	Cons            string
	CustomRules     string
	// OrderStatus — статус заказа из WB: "buyout" (выкуп), "returned" (товар возвращён). Если "returned" — не предлагать оформлять возврат.
	OrderStatus     string
}

// FeedbackReplyResult — результат генерации ответа (текст + метаданные).
type FeedbackReplyResult struct {
	Text       string
	TokensUsed int
	ModelUsed  string
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
