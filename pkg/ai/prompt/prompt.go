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
{{.order_status_hint}}

ОБЩИЙ СТИЛЬ:
Пиши как живой человек — коротко, тепло, без канцелярита.
Избегай: "максимально", "данный", "являемся", "надеемся что", двойных "спасибо" в одном ответе, общих фраз которые подошли бы любому товару.
Не пересказывай отзыв покупателю обратно.

ПОЗИТИВ (4–5 звёзд):
Поблагодари живо и естественно. Зацепись за конкретную деталь или эмоцию из отзыва — покажи что прочитал именно его отзыв, а не шаблон.
Не хвали товар сам — покупатель уже это сделал.
Не предлагай дополнительные товары.

НЕГАТИВ (1–3 звезды):
Не отправляй в поддержку. 
Признай конкретную проблему из отзыва — не абстрактно "нам жаль", 
а именно то, на что жалуется покупатель.
Извинись один раз, искренне, без "приносим свои извинения" и подобного канцелярита.
Если проблема связана с качеством товара — предложи проверить на брак или оформить возврат через маркетплейс (только если товар не возвращён: см. данные отзыва).
Если покупатель недоволен чем-то другим (доставка, упаковка, несоответствие ожиданиям) — 
признай это, не переваливай вину на маркетплейс.
Тон: спокойный, человечный, без заискивания.

Доп. правила бренда: {{.custom_rules}}.

Требования: русский язык, до 400 символов, без ссылок, без телефонов, без email, без просьб перейти в мессенджеры.`

var feedbackReplyPrompt = prompts.NewPromptTemplate(
	feedbackReplyTemplate,
	[]string{
		"brand_name", "tone_instruction", "product_name", "rating",
		"text", "pros", "cons", "custom_rules", "order_status_hint",
	},
)

// orderStatusHint возвращает подсказку для промпта, если товар уже возвращён — чтобы не предлагать возврат в ответе.
func orderStatusHint(orderStatus string) string {
	if orderStatus == "returned" {
		return "ВАЖНО: по данным маркетплейса товар уже возвращён покупателем (отказ от товара). Не предлагай оформлять возврат через маркетплейс — возвращать нечего. Признай сожаление и поблагодари за отзыв."
	}
	return ""
}

// FormatFeedbackReply подставляет переменные в шаблон ответа на отзыв.
func FormatFeedbackReply(in *ai.FeedbackReplyInput) (string, error) {
	if in == nil {
		return "", fmt.Errorf("feedback reply input is nil")
	}
	return feedbackReplyPrompt.Format(map[string]any{
		"brand_name":        in.BrandName,
		"tone_instruction":  in.ToneInstruction,
		"product_name":     in.ProductName,
		"rating":            in.Rating,
		"text":              in.Text,
		"pros":              in.Pros,
		"cons":              in.Cons,
		"custom_rules":      in.CustomRules,
		"order_status_hint": orderStatusHint(in.OrderStatus),
	})
}
