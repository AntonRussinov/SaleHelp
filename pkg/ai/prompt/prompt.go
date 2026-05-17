package prompt

import (
	"fmt"
	"strings"

	"salehelp/pkg/ai"
)

// basePromptTemplate — общая структура промпта; плейсхолдеры (по порядку):
// brand_block, tone_instruction, addressing_instruction, emoji_instruction,
// product_name, review_data, category_instructions, requirements_block, custom_rules_section.
const basePromptTemplate = `Ты — представитель бренда.
%s
Стиль: %s.
Обращение к покупателю: %s.
Эмодзи: %s.
Товар: %s.

# Данные отзыва
%s

%s
# Требования к ответу
- %s
- Русский язык
- Без ссылок, телефонов, email, призывов перейти в мессенджеры
- Пиши как живой человек: коротко, тепло, без канцелярита
- Не хвали товар сам — покупатель уже это сделал
- Не пересказывай отзыв обратно покупателю
- Избегай: «максимально», «данный», «являемся», «надеемся что», двойного «спасибо» в одном ответе%s%s

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

func buildBrandBlock(brandName, brandDescription string) string {
	if brandName == "" {
		brandName = "магазина"
	}
	if strings.TrimSpace(brandDescription) == "" {
		return "Бренд: " + brandName + "."
	}
	return fmt.Sprintf("Бренд: %s. %s", brandName, strings.TrimSpace(brandDescription))
}

func buildAddressingInstruction(addressing string) string {
	if addressing == "ty" {
		return "обращайся на «ты»"
	}
	return "обращайся на «вы»"
}

func buildEmojiInstruction(emojiMode string) string {
	switch emojiMode {
	case "none":
		return "без эмодзи"
	case "expressive":
		return "2–4 уместных эмодзи"
	default: // minimal
		return "не более 1–2 уместных эмодзи"
	}
}

func buildResponseLengthInstruction(responseLength string) string {
	switch responseLength {
	case "short":
		return "длина ответа: до 200 символов"
	case "long":
		return "длина ответа: до 600 символов"
	default: // medium
		return "длина ответа: до 400 символов"
	}
}

// buildCategoryInstructions возвращает инструкции для нужного типа отзыва + кастомный промпт категории.
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

// buildCustomRulesSection формирует секцию «Правила бренда» и «Запрещено упоминать».
func buildCustomRulesSection(rules, forbidden []string) string {
	var sb strings.Builder

	var ruleItems []string
	for _, r := range rules {
		if t := strings.TrimSpace(r); t != "" {
			ruleItems = append(ruleItems, "- "+t)
		}
	}
	if len(ruleItems) > 0 {
		sb.WriteString("\n\n# Правила бренда\n")
		sb.WriteString(strings.Join(ruleItems, "\n"))
	}

	var forbItems []string
	for _, f := range forbidden {
		if t := strings.TrimSpace(f); t != "" {
			forbItems = append(forbItems, "- "+t)
		}
	}
	if len(forbItems) > 0 {
		sb.WriteString("\n\n# Запрещено упоминать\n")
		sb.WriteString(strings.Join(forbItems, "\n"))
	}

	return sb.String()
}

// buildSignatureInstruction возвращает инструкцию подписи если задана.
func buildSignatureInstruction(signature string) string {
	if t := strings.TrimSpace(signature); t != "" {
		return "\n- Завершай ответ подписью: " + t
	}
	return ""
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
	brandBlock := buildBrandBlock(in.BrandName, in.BrandDescription)
	addressingInstruction := buildAddressingInstruction(in.Addressing)
	emojiInstruction := buildEmojiInstruction(in.EmojiMode)
	lengthInstruction := buildResponseLengthInstruction(in.ResponseLength)
	reviewData := buildReviewData(in.Rating, in.Text, in.Pros, in.Cons, in.OrderStatus)
	categorySection := buildCategoryInstructions(in.Rating, in.CategoryCustomPrompt)
	customRulesSection := buildCustomRulesSection(in.CustomRules, in.ForbiddenWords)
	signatureInstruction := buildSignatureInstruction(in.Signature)

	return fmt.Sprintf(
		basePromptTemplate,
		brandBlock,
		in.ToneInstruction,
		addressingInstruction,
		emojiInstruction,
		in.ProductName,
		reviewData,
		categorySection,
		lengthInstruction,
		customRulesSection,
		signatureInstruction,
	), nil
}
