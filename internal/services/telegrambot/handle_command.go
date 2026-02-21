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
	case len(cmd) > 1 && cmd[0] == "/start" && strings.HasPrefix(cmd[1], "invite_"):
		token := strings.TrimPrefix(cmd[1], "invite_")
		s.handleJoin(ctx, chatID, token)
		return

	case len(cmd) > 0 && cmd[0] == "/start":
		_ = s.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			"Привет! 👋\n\nВыбери действие:",
			mainMenuKeyboard(),
		)
		return

	case len(cmd) > 0 && cmd[0] == "/cancel":
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			"Ок, отменил. Можешь начать заново.",
			mainMenuKeyboard(),
		)
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

	case len(cmd) == 2 && cmd[0] == "/join":
		s.handleJoin(ctx, chatID, cmd[1])
		return

	case len(cmd) == 3 && cmd[0] == "/shop":
		s.handleShopToggle(ctx, chatID, cmd[1], cmd[2])
		return

	case len(cmd) == 3 && cmd[0] == "/notify":
		s.handleNotifyToggle(ctx, chatID, cmd[1], cmd[2])
		return

	case len(cmd) == 2 && cmd[0] == "/chats":
		s.handleChats(ctx, chatID, cmd[1])
		return

	case len(cmd) >= 1 && cmd[0] == "/abc_now":
		s.handleABCNow(ctx, chatID, cmd)
		return

	case len(cmd) > 0 && cmd[0] == "/admin_stats":
		s.handleAdminStats(ctx, chatID)
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
