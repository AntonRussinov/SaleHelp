package telegram

import (
	"context"
	"fmt"

	"salehelp/pkg/notify"
)

const telegramMaxLen = 4096

type Sender struct {
	client *Client
}

func NewSender(token string) (*Sender, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram token is empty")
	}
	return &Sender{
		client: NewClient(token),
	}, nil
}

// DefaultRussianCommands возвращает команды для Bot Command Menu (☰) на русском.
func DefaultRussianCommands() []BotCommand {
	return []BotCommand{
		{Command: "start", Description: "Старт"},
		{Command: "addshop", Description: "Добавить магазин"},
		{Command: "myshops", Description: "Мои магазины"},
		{Command: "orders", Description: "Текущие заказы на сборке"},
		{Command: "alerts", Description: "Настроить уведомления"},
		{Command: "invite", Description: "Пригласить в магазин"},
		{Command: "abc_now", Description: "ABC-анализ"},
		{Command: "cancel", Description: "Отмена"},
	}
}

// SetMyCommands задаёт список команд в меню бота (иконка ☰ рядом с полем ввода).
func (s *Sender) SetMyCommands(ctx context.Context, commands []BotCommand) error {
	return s.client.SetMyCommands(ctx, commands)
}

// SendWithInlineKeyboard отправляет сообщение с inline-кнопками (callback_data при нажатии).
// inlineKeyboard — reply_markup, например: map[string]any{"inline_keyboard": [][]map[string]string{...}}
func (s *Sender) SendWithInlineKeyboard(
	ctx context.Context,
	chatID int64,
	text string,
	inlineKeyboard map[string]any,
) error {
	if chatID == 0 {
		return fmt.Errorf("telegram chat_id is 0")
	}
	if text == "" {
		return fmt.Errorf("telegram text is empty")
	}
	text = escapeHTML(text)
	payload := map[string]any{
		"chat_id":                  chatID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
		"reply_markup":             inlineKeyboard,
	}
	return s.client.SendRaw(ctx, payload)
}

// AnswerCallbackQuery снимает индикатор загрузки с inline-кнопки.
func (s *Sender) AnswerCallbackQuery(ctx context.Context, callbackQueryID string) error {
	return s.client.AnswerCallbackQuery(ctx, callbackQueryID)
}

// EditMessageReplyMarkup заменяет inline-клавиатуру уже отправленного сообщения.
func (s *Sender) EditMessageReplyMarkup(ctx context.Context, chatID int64, messageID int, replyMarkup any) error {
	return s.client.EditMessageReplyMarkup(ctx, chatID, messageID, replyMarkup)
}

func (s *Sender) Channel() notify.Channel {
	return notify.ChannelTelegram
}

func (s *Sender) Send(ctx context.Context, r notify.Recipient, m notify.Message) error {
	if r.TelegramChatID == nil {
		return fmt.Errorf("telegram recipient chat_id is nil")
	}

	// делим ПЛЕЙН-ТЕКСТ, чтобы не резать HTML
	parts := splitTextByRuneLimit(m.Text, telegramMaxLen-200) // запас под title + html escape

	for i, part := range parts {
		text := escapeHTML(part)

		// Title показываем только в первой части (чтобы не спамить жирным заголовком)
		if i == 0 && m.Title != "" {
			text = fmt.Sprintf("<b>%s</b>\n\n%s", escapeHTML(m.Title), text)
		}

		// доп. страховка: если после escape всё равно больше лимита — режем уже готовый текст
		if len([]rune(text)) > telegramMaxLen {
			rs := []rune(text)
			text = string(rs[:telegramMaxLen-1]) + "…"
		}

		isLastPart := i == len(parts)-1
		if isLastPart && m.InlineKeyboard != nil {
			if err := s.client.SendRaw(ctx, map[string]any{
				"chat_id":                  *r.TelegramChatID,
				"text":                     text,
				"parse_mode":               "HTML",
				"disable_web_page_preview": true,
				"reply_markup":             m.InlineKeyboard,
			}); err != nil {
				return err
			}
		} else if err := s.client.SendMessage(
			ctx,
			*r.TelegramChatID,
			text,
			"HTML",
			true,
		); err != nil {
			return err
		}
	}

	return nil
}

func splitTextByRuneLimit(s string, limit int) []string {
	if limit <= 0 {
		return []string{s}
	}
	rs := []rune(s)
	if len(rs) <= limit {
		return []string{s}
	}

	out := make([]string, 0, (len(rs)/limit)+1)
	for len(rs) > 0 {
		n := limit
		if len(rs) < n {
			n = len(rs)
		}
		out = append(out, string(rs[:n]))
		rs = rs[n:]
	}
	return out
}

func (s *Sender) SendWithKeyboard(
	ctx context.Context,
	chatID int64,
	text string,
	keyboard map[string]any,
) error {
	if chatID == 0 {
		return fmt.Errorf("telegram chat_id is 0")
	}
	if text == "" {
		return fmt.Errorf("telegram text is empty")
	}

	// экранируем так же, как в обычной Send
	text = escapeHTML(text)

	payload := map[string]any{
		"chat_id":                  chatID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
		"reply_markup":             keyboard,
	}

	return s.client.SendRaw(ctx, payload)
}

// Минимальный HTML-эскейп (достаточно для <, >, &).
func escapeHTML(s string) string {
	// не используем html.EscapeString? можно, но это минималка без лишних импортов
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch r {
		case '&':
			out = append(out, []rune("&amp;")...)
		case '<':
			out = append(out, []rune("&lt;")...)
		case '>':
			out = append(out, []rune("&gt;")...)
		default:
			out = append(out, r)
		}
	}
	return string(out)
}

// func splitByRuneLimit(s string, limit int) []string {
// 	if limit <= 0 {
// 		return []string{s}
// 	}
// 	rs := []rune(s)
// 	if len(rs) <= limit {
// 		return []string{s}
// 	}

// 	out := make([]string, 0, (len(rs)/limit)+1)
// 	for len(rs) > 0 {
// 		n := limit
// 		if len(rs) < n {
// 			n = len(rs)
// 		}
// 		out = append(out, string(rs[:n]))
// 		rs = rs[n:]
// 	}
// 	return out
// }
