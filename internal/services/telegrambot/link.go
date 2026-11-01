package telegrambot

import (
	"context"
	"strings"
	"time"

	"salehelp/domain"
)

// handleLink обрабатывает /link <код>. Код получают на сайте (Настройки → Привязать Telegram).
// Если текущий чат уже привязан к другому пользователю (telegram-only), выполняется слияние:
// магазины переносятся на целевого пользователя, чат привязывается к нему, у старого пользователя чат отвязывается.
func (s *Service) handleLink(ctx context.Context, chatID int64, code string) {
	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		_ = s.tg.SendMessage(ctx, chatID, "Укажи код: /link КОД (код можно получить на сайте в настройках).")
		return
	}

	if s.linkTokensRepo == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Привязка временно недоступна.")
		return
	}

	targetUserID, err := s.linkTokensRepo.Consume(ctx, code, time.Now())
	if err != nil {
		s.log.Error("link token consume failed", "chat_id", chatID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка проверки кода. Попробуй ещё раз.")
		return
	}
	if targetUserID <= 0 {
		_ = s.tg.SendMessage(ctx, chatID, "Код не найден или истёк. Получи новый код на сайте и введи его снова.")
		return
	}

	currentUser, err := s.usersRepo.GetByTelegramChatID(ctx, chatID)
	if err != nil {
		s.log.Error("get user by telegram chat id failed", "chat_id", chatID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка. Попробуй позже.")
		return
	}

	if currentUser != nil && currentUser.ID == targetUserID {
		_ = s.tg.SendMessage(ctx, chatID, "Этот чат уже привязан к твоему аккаунту.")
		return
	}

	if currentUser != nil && currentUser.ID != targetUserID {
		// Слияние: переносим магазины, отвязываем чат от старого пользователя, привязываем к целевому
		if err := s.shopsRepo.ReassignShopsToUser(ctx, currentUser.ID, targetUserID); err != nil {
			s.log.Error("reassign shops failed", "from", currentUser.ID, "to", targetUserID, "err", err)
			_ = s.tg.SendMessage(ctx, chatID, "Ошибка объединения аккаунтов. Обратись в поддержку.")
			return
		}
		if err := s.usersRepo.ClearTelegramChatID(ctx, currentUser.ID); err != nil {
			s.log.Warn("clear old user telegram_chat_id failed", "user_id", currentUser.ID, "err", err)
		}
		if err := s.usersRepo.SetTelegramChatID(ctx, targetUserID, chatID); err != nil {
			s.log.Error("set telegram_chat_id failed", "user_id", targetUserID, "err", err)
			_ = s.tg.SendMessage(ctx, chatID, "Ошибка привязки. Попробуй ещё раз.")
			return
		}
		s.ensureOwnerEndpointsForUser(ctx, targetUserID, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Готово ✅ Аккаунт объединён: магазины и настройки теперь в одном аккаунте. Можешь пользоваться ботом и сайтом с одними данными.")
		return
	}

	// Текущий чат не привязан ни к кому — просто привязываем к целевому пользователю
	if err := s.usersRepo.SetTelegramChatID(ctx, targetUserID, chatID); err != nil {
		s.log.Error("set telegram_chat_id failed", "user_id", targetUserID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка привязки. Попробуй ещё раз.")
		return
	}
	s.ensureOwnerEndpointsForUser(ctx, targetUserID, chatID)
	_ = s.tg.SendMessage(ctx, chatID, "Готово ✅ Telegram привязан к твоему аккаунту. Настройки бота и сайта теперь общие.")
}

// startLinkEmail запускает привязку email: либо сохраняем email за текущим аккаунтом, либо (если email занят) выдаём код для подтверждения на сайте.
func (s *Service) startLinkEmail(ctx context.Context, chatID int64) {
	st := &domain.TelegramChatState{
		ChatID:  chatID,
		State:   StateWaitEmailForLink,
		Payload: []byte(`{}`),
	}
	if err := s.statesRepo.Upsert(ctx, st, time.Now()); err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка. Попробуй ещё раз.")
		return
	}
	_ = s.tg.SendMessage(ctx, chatID,
		"Введи email от аккаунта на сайте.\n\n"+
			"• Если такого аккаунта ещё нет — мы сохраним email за тобой.\n"+
			"• Если аккаунт с этим email уже есть — мы вышлем тебе код; зайди на сайт под этим email и введи код в настройках (подтверждение, что это твой аккаунт).\n\n"+
			"Отмена: /cancel")
}

// stepEmailForLink обрабатывает ввод email в состоянии StateWaitEmailForLink.
// Мы не делаем слияние по вводу чужого email (это небезопасно: кто угодно мог бы привязать свой чат к чужому аккаунту).
// Если email уже занят — предлагаем получить код на сайте и использовать /link КОД.
func (s *Service) stepEmailForLink(ctx context.Context, chatID int64, text string) {
	email := strings.TrimSpace(strings.ToLower(text))
	if email == "" || !strings.Contains(email, "@") {
		_ = s.tg.SendMessage(ctx, chatID, "Введи корректный email или /cancel.")
		return
	}

	currentUser, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil {
		s.log.Error("ensure user failed", "chat_id", chatID, "err", err)
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка. Попробуй /link_email снова.")
		return
	}

	webUser, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("get user by email failed", "email", email, "err", err)
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка. Попробуй ещё раз или /cancel.")
		return
	}

	_ = s.statesRepo.Clear(ctx, chatID)

	if webUser != nil && webUser.ID != currentUser.ID {
		// Email занят — выдаём код; подтверждение на сайте (владелец email вводит код после входа).
		if s.linkTokensRepo == nil {
			_ = s.tg.SendMessage(ctx, chatID, "Привязка по коду временно недоступна. Зайди на сайт → Настройки → Привязать Telegram → получи код и введи здесь: /link КОД")
			return
		}
		code, _, err := s.linkTokensRepo.CreateFromTelegram(ctx, chatID)
		if err != nil {
			s.log.Error("create link token from telegram failed", "chat_id", chatID, "err", err)
			_ = s.tg.SendMessage(ctx, chatID, "Не удалось сгенерировать код. Попробуй позже или используй способ с сайта: зайди на сайт → Привязать Telegram → /link КОД")
			return
		}
		_ = s.tg.SendMessage(ctx, chatID,
			"Код для привязки (действует ~10 мин):\n\n"+
				code+"\n\n"+
				"Зайди на сайт под email "+email+", открой Настройки → Подтвердить привязку Telegram и введи этот код. Тогда этот чат привяжется к твоему аккаунту.")
		return
	}

	if webUser != nil && webUser.ID == currentUser.ID {
		_ = s.tg.SendMessage(ctx, chatID, "Этот чат уже привязан к аккаунту с таким email.")
		return
	}

	// Email свободен — привязываем к текущему (telegram) пользователю
	currentUser.Email = &email
	if err := s.usersRepo.Update(ctx, currentUser); err != nil {
		s.log.Error("update user email failed", "user_id", currentUser.ID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Не получилось сохранить email. Попробуй позже.")
		return
	}
	_ = s.tg.SendMessage(ctx, chatID, "Email привязан к твоему аккаунту. Можешь войти на сайт с этим email (создай пароль через «Восстановить пароль», если ещё не регистрировался).")
}
