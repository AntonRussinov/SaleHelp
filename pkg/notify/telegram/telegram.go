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

		if err := s.client.SendMessage(
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
