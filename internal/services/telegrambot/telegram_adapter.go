package telegrambot

import (
	"context"

	"salehelp/pkg/notify"
	"salehelp/pkg/notify/telegram"
)

type TelegramAdapter struct {
	sender *telegram.Sender
}

func NewTelegramAdapter(sender *telegram.Sender) *TelegramAdapter {
	return &TelegramAdapter{sender: sender}
}

// Реализует TelegramAPI интерфейс
func (a *TelegramAdapter) SendMessage(ctx context.Context, chatID int64, text string) error {
	r := notify.Recipient{
		TelegramChatID: &chatID,
	}

	m := notify.Message{
		Text: text,
	}

	return a.sender.Send(ctx, r, m)
}
func (a *TelegramAdapter) SendMessageWithKeyboard(
	ctx context.Context,
	chatID int64,
	text string,
	keyboard map[string]any,
) error {
	return a.sender.SendWithKeyboard(ctx, chatID, text, keyboard)
}
