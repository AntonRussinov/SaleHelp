package prompt

import (
	"fmt"
	"strings"

	"salehelp/pkg/ai"
)

// basePromptTemplate — общая структура промпта; плейсхолдеры (по порядку):
// brand_name, tone_instruction, product_name, review_data, category_instructions, custom_rules_section.
const basePromptTemplate = `Ты — менеджер бренда %s.
Стиль: %s.
Товар: %s.

# Данные отзыва
%s

%s
# Общие требования к ответу
- Русский язык, до 400 символов
- Без ссылок, телефонов, email, призывов перейти в мессенджеры
- Пиши как живой человек: коротко, тепло, без канцелярита
- Не хвали товар сам — покупатель уже это сделал
- Не пересказывай отзыв обратно покупателю
- Избегай: «максимально», «данный», «являемся», «надеемся что», двойного «спасибо» в одном ответе%s

# Ответ:`

const positiveInstructions = `# Как отвечать (позитивный отзыв, 4–5 звёзд)
- Поблагодари живо и естественно, без шаблонных заготовок
- Не предлагай дополнительные товары`

const neutralInstructions = `# Как отвечать (нейтральный отзыв, 3 звезды)
- Поблагодари за честный отзыв
- Если упомянут конкретный недостаток — признай его, не абстрактно
- Держи уважительный тон без заискивания`

const negativeInstructions = `# Как отвечать (негативный отзыв, 1–2 звезды)
- Признай конкретную проблему из отзыва — именно то, на что жалуется покупатель, не абстрактно «нам жаль»
- Извинись один раз, искренне, без канцелярита (не «приносим свои извинения»)
- Не отправляй в поддержку и не предлагай позвонить
- Если проблема — качество товара: скажи, что постараемся исправить; не предлагай возврат
- Если проблема — доставка, упаковка, несоответствие ожиданиям: признай это, не перекладывай вину на маркетплейс
- Тон: спокойный, человечный`

// buildCategoryInstructions возвращает инструкции для нужного типа отзыва + опциональный кастомный промпт категории.
func buildCategoryInstructions(rating int, categoryCustomPrompt string) string {
	var base string
	switch {
	case rating >= 4:
		base = positiveInstructions
	case rating == 3:
		base = neutralInstructions
	default:
		base = negativeInstructions
	}
	if t := strings.TrimSpace(categoryCustomPrompt); t != "" {
		return base + "\n" + t
	}
	return base
}

// buildCustomRulesSection формирует секцию «Правила бренда» как маркированный список.
func buildCustomRulesSection(rules []string) string {
	var items []string
	for _, r := range rules {
		if t := strings.TrimSpace(r); t != "" {
			items = append(items, "- "+t)
		}
	}
	if len(items) == 0 {
		return ""
	}
	return "\n\n# Правила бренда\n" + strings.Join(items, "\n")
}

// buildReviewData формирует раздел с данными отзыва; включает только непустые поля.
func buildReviewData(rating int, text, pros, cons, orderStatus string) string {
	parts := []string{fmt.Sprintf("Оценка: %d", rating)}
	if t := strings.TrimSpace(text); t != "" {
		parts = append(parts, "Комментарий: "+t)
	}
	if t := strings.TrimSpace(pros); t != "" {
		parts = append(parts, "Плюсы: "+t)
	}
	if t := strings.TrimSpace(cons); t != "" {
		parts = append(parts, "Минусы: "+t)
	}

	result := strings.Join(parts, "\n")

	if orderStatus == "returned" {
		result += "\n\nВАЖНО: по данным маркетплейса товар уже возвращён покупателем. Не предлагай возврат — он уже оформлен. Признай сожаление и поблагодари за отзыв."
	}
	return result
}

// FormatFeedbackReply подставляет данные в промпт-шаблон ответа на отзыв.
func FormatFeedbackReply(in *ai.FeedbackReplyInput) (string, error) {
	if in == nil {
		return "", fmt.Errorf("feedback reply input is nil")
	}
	reviewData := buildReviewData(in.Rating, in.Text, in.Pros, in.Cons, in.OrderStatus)
	categorySection := buildCategoryInstructions(in.Rating, in.CategoryCustomPrompt)
	customRulesSection := buildCustomRulesSection(in.CustomRules)

	return fmt.Sprintf(
		basePromptTemplate,
		in.BrandName,
		in.ToneInstruction,
		in.ProductName,
		reviewData,
		categorySection,
		customRulesSection,
	), nil
}
