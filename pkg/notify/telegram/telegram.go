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

	text := m.Text
	if m.Title != "" {
		// HTML-friendly формат (без Markdown экранирования)
		text = fmt.Sprintf("<b>%s</b>\n\n%s", escapeHTML(m.Title), escapeHTML(text))
	} else {
		text = escapeHTML(text)
	}

	// Telegram limit: 4096 символов
	if len([]rune(text)) > telegramMaxLen {
		rs := []rune(text)
		text = string(rs[:telegramMaxLen-1]) + "…"
	}

	return s.client.SendMessage(
		ctx,
		*r.TelegramChatID,
		text,
		"HTML", // parse mode
		true,   // disable web page preview
	)
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
