package telegrambot

import (
	"context"
	"salehelp/domain"
	"strconv"
	"time"
)

func (s *Service) handleJoin(
	ctx context.Context,
	chatID int64,
	token string,
) {
	invite, err := s.invitesRepo.GetByToken(ctx, token)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка проверки токена.")
		return
	}

	if invite == nil || invite.ExpiresAt.Before(time.Now()) {
		_ = s.tg.SendMessage(ctx, chatID, "Токен недействителен или истёк.")
		return
	}

	// создаём/получаем пользователя
	if _, err := s.ensureUserByChatID(ctx, chatID); err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка пользователя.")
		return
	}

	ok, msg := s.canInviteMoreTelegramChats(ctx, invite.ShopID) // тот же helper
	if !ok {
		_ = s.tg.SendMessage(ctx, chatID, msg)
		return
	}

	// проверяем, не подключён ли уже этот чат
	exists, err := s.endpointsRepo.Exists(
		ctx,
		invite.ShopID,
		"telegram",
		strconv.FormatInt(chatID, 10),
	)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка подключения.")
		return
	}
	if exists {
		_ = s.tg.SendMessage(ctx, chatID, "Этот чат уже подключён к магазину.")
		return
	}

	// создаём endpoint
	err = s.endpointsRepo.Create(ctx, &domain.NotificationEndpoint{
		ShopID:  invite.ShopID,
		Channel: "telegram",
		Address: strconv.FormatInt(chatID, 10),
		Enabled: true,
	})
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка сохранения подключения.")
		return
	}

	_ = s.tg.SendMessageWithKeyboard(ctx, chatID,
		"Готово ✅\n\n"+
			"Ты подключён к магазину и будешь получать уведомления.",
		mainMenuKeyboard(),
	)
}
