package telegrambot

import (
	"context"
	"salehelp/internal/handler/telegram/dto"
	"strings"
)

func (s *Service) HandleUpdate(ctx context.Context, upd *dto.Update) {
	if upd == nil {
		return
	}

	// Нажатие на inline-кнопку (например выбор магазина для ABC)
	if upd.CallbackQuery != nil {
		s.handleCallbackQuery(ctx, upd.CallbackQuery)
		return
	}

	if upd.Message == nil {
		return
	}

	chatID := upd.Message.Chat.ID
	text := strings.TrimSpace(upd.Message.Text)
	if text == "" {
		return
	}

	// команды
	if strings.HasPrefix(text, "/") {
		s.handleCommand(ctx, chatID, text)
		return
	}

	// wizard
	s.handleState(ctx, chatID, text)
}
