package telegrambot

import (
	"context"
	"salehelp/domain"
	"strings"
	"time"
)

func (s *Service) handleCommand(ctx context.Context, chatID int64, text string) {
	cmd := strings.Fields(text)
	switch {
	case len(cmd) > 0 && cmd[0] == "/start":
		_ = s.tg.SendMessage(ctx, chatID,
			"Привет! 👋\n\nКоманды:\n"+
				"/addshop — добавить магазин\n"+
				"/cancel — отменить текущий шаг",
		)
		return

	case len(cmd) > 0 && cmd[0] == "/cancel":
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Ок, отменил. Можешь начать заново: /addshop")
		return

	case len(cmd) > 0 && cmd[0] == "/addshop":
		s.startAddShop(ctx, chatID)
		return

	case len(cmd) > 0 && cmd[0] == "/myshops":
		_ = s.handleMyShops(ctx, chatID)
		return

	case len(cmd) > 0 && cmd[0] == "/invite":
		s.startInvite(ctx, chatID)
		return

	default:
		_ = s.tg.SendMessage(ctx, chatID, "Неизвестная команда. Напиши /start")
	}
}

func (s *Service) startAddShop(ctx context.Context, chatID int64) {
	st := &domain.TelegramChatState{
		ChatID:  chatID,
		State:   StateWaitWBToken,
		Payload: []byte(`{}`),
	}

	if err := s.statesRepo.Upsert(ctx, st, time.Now()); err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Не получилось начать регистрацию. Попробуй ещё раз.")
		return
	}

	_ = s.tg.SendMessage(ctx, chatID,
		"Давай подключим магазин.\n\n"+
			"Шаг 1/2: отправь WB API token (как текстовым сообщением).",
	)
}
