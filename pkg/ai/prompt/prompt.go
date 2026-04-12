package prompt

import (
	"fmt"
	"strings"

	"salehelp/pkg/ai"

	"github.com/tmc/langchaingo/prompts"
)

const feedbackReplyTemplate = `Ты — менеджер бренда {{.brand_name}}.
Твой стиль: {{.tone_instruction}}.
Твоя задача — ответить на отзыв покупателя о товаре {{.product_name}}.

{{.review_data}}
{{.order_status_hint}}

ОБЩИЙ СТИЛЬ:
Пиши как живой человек — коротко, тепло, без канцелярита. Не хвали сам себя. 
Избегай: "максимально", "данный", "являемся", "надеемся что", двойных "спасибо" в одном ответе, общих фраз которые подошли бы любому товару.
Не пересказывай отзыв покупателю обратно.

ПОЗИТИВ (4–5 звёзд):
Поблагодари живо и естественно, это не должно быть слишком шаблонно, но и не надо слишком выдумывать.
Не хвали товар сам — покупатель уже это сделал.
Не предлагай дополнительные товары.

НЕГАТИВ (1–3 звезды):
Не отправляй в поддержку. 
Признай конкретную проблему из отзыва — не абстрактно "нам жаль", 
а именно то, на что жалуется покупатель.
Извинись один раз, искренне, без "приносим свои извинения" и подобного канцелярита.
Если проблема связана с качеством товара — не надо предлагать возврат товара, скажи, что мы постараемся это исправить или улучшить.
Если покупатель недоволен чем-то другим (доставка, упаковка, несоответствие ожиданиям) — 
признай это, не переваливай вину на маркетплейс.
Тон: спокойный, человечный, без заискивания.

Доп. правила бренда: {{.custom_rules}}.

Требования: русский язык, до 400 символов, без ссылок, без телефонов, без email, без просьб перейти в мессенджеры.`

var feedbackReplyPrompt = prompts.NewPromptTemplate(
	feedbackReplyTemplate,
	[]string{
		"brand_name", "tone_instruction", "product_name", "review_data",
		"custom_rules", "order_status_hint",
	},
)

// buildReviewData формирует строку «Данные отзыва», включая только непустые поля (оценка всегда, текст/плюсы/минусы — при наличии).
func buildReviewData(rating int, text, pros, cons string) string {
	parts := []string{fmt.Sprintf("Оценка: %d", rating)}
	if strings.TrimSpace(text) != "" {
		parts = append(parts, "Текст: "+strings.TrimSpace(text))
	}
	if strings.TrimSpace(pros) != "" {
		parts = append(parts, "Плюсы: "+strings.TrimSpace(pros))
	}
	if strings.TrimSpace(cons) != "" {
		parts = append(parts, "Минусы: "+strings.TrimSpace(cons))
	}
	return "Данные отзыва: " + strings.Join(parts, ". ")
}

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
		"product_name":      in.ProductName,
		"review_data":       buildReviewData(in.Rating, in.Text, in.Pros, in.Cons),
		"custom_rules":      in.CustomRules,
		"order_status_hint": orderStatusHint(in.OrderStatus),
	})
}
