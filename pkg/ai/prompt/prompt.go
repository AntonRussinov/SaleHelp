package prompt

import (
	"fmt"

	"salehelp/pkg/ai"

	"github.com/tmc/langchaingo/prompts"
)

const feedbackReplyTemplate = `Ты — менеджер бренда {{.brand_name}}.
Твой стиль: {{.tone_instruction}}.
Твоя задача — ответить на отзыв покупателя о товаре {{.product_name}}.
Данные отзыва: Оценка: {{.rating}}. Текст: {{.text}}. Плюсы: {{.pros}}. Минусы: {{.cons}}.
Негатив (1–3 звезды): не отправляй в поддержку. Признай проблему, извинись и предложи проверить товар на брак или оформить возврат через интерфейс маркетплейса.
Позитив (4–5 звёзд): поблагодари и подчеркни одну характеристику товара. Не предлагай дополнительные товары.
Доп. правила бренда: {{.custom_rules}}.
Требования: русский язык, до 400 символов, без ссылок, без телефонов, без email, без просьб перейти в мессенджеры.`

var feedbackReplyPrompt = prompts.NewPromptTemplate(
	feedbackReplyTemplate,
	[]string{
		"brand_name", "tone_instruction", "product_name", "rating",
		"text", "pros", "cons", "custom_rules",
	},
)

// FormatFeedbackReply подставляет переменные в шаблон ответа на отзыв.
func FormatFeedbackReply(in *ai.FeedbackReplyInput) (string, error) {
	if in == nil {
		return "", fmt.Errorf("feedback reply input is nil")
	}
	return feedbackReplyPrompt.Format(map[string]any{
		"brand_name":       in.BrandName,
		"tone_instruction": in.ToneInstruction,
		"product_name":     in.ProductName,
		"rating":           in.Rating,
		"text":             in.Text,
		"pros":             in.Pros,
		"cons":             in.Cons,
		"custom_rules":     in.CustomRules,
	})
}
