package telegrambot

import (
	"context"
	"encoding/json"
	"fmt"
	"salehelp/domain"
	"strings"
	"time"
)

const (
	StateIdle             = "idle"
	StateWaitWBToken      = "wait_wb_token"
	StateWaitShopName     = "wait_shop_name"
	StateInviteChooseShop = "invite_choose_shop"
	StateWaitEmailForLink = "wait_email_for_link"
)

const (
	PayloadWBToken = "wb_token"
)

type payloadMap map[string]string

func (s *Service) handleState(ctx context.Context, chatID int64, text string) {
	st, err := s.statesRepo.GetByChatID(ctx, chatID)
	if err != nil {
		s.log.Error("get chat state failed", "chat_id", chatID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка. Напиши /addshop чтобы начать заново.")
		return
	}

	if st == nil || st.State == "" || st.State == StateIdle {
		_ = s.tg.SendMessage(ctx, chatID, "Напиши /addshop чтобы подключить магазин.")
		return
	}

	switch st.State {
	case StateWaitWBToken:
		s.stepWBToken(ctx, chatID, st, text)
		return
	case StateWaitShopName:
		s.stepShopName(ctx, chatID, st, text)
		return
	case StateInviteChooseShop:
		s.handleInviteChooseShop(ctx, chatID, text, st)
		return
	case StateWaitEmailForLink:
		s.stepEmailForLink(ctx, chatID, text)
		return
	default:
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Состояние сбилось. Напиши /addshop чтобы начать заново.")
		return
	}
}

func (s *Service) stepWBToken(ctx context.Context, chatID int64, st *domain.TelegramChatState, token string) {
	token = strings.TrimSpace(token)
	if len(token) < 20 {
		_ = s.tg.SendMessage(ctx, chatID, "Похоже, это не WB token. Пришли токен ещё раз или напиши /cancel.")
		return
	}

	p, err := parsePayload(st.Payload)
	if err != nil {
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка payload. Напиши /addshop чтобы начать заново.")
		return
	}

	if err := s.monitoring.CheckWBToken(ctx, token); err != nil {
		_ = s.tg.SendMessage(ctx, chatID,
			"❌ Токен не подходит.\n\n"+
				"Проверь:\n"+
				"• что токен скопирован полностью\n"+
				"• что он из кабинета WB\n"+
				"• что у него есть доступ к карточкам",
		)
		return
	}

	p[PayloadWBToken] = token

	next := &domain.TelegramChatState{
		ChatID:  chatID,
		State:   StateWaitShopName,
		Payload: mustJSON(p),
	}

	if err := s.statesRepo.Upsert(ctx, next, time.Now()); err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Не получилось сохранить токен. Попробуй ещё раз.")
		return
	}

	_ = s.tg.SendMessage(ctx, chatID, "Шаг 2/2: отправь название магазина (любое, как тебе удобно).")
}

func (s *Service) stepShopName(ctx context.Context, chatID int64, st *domain.TelegramChatState, name string) {
	name = strings.TrimSpace(name)
	if len([]rune(name)) < 2 {
		_ = s.tg.SendMessage(ctx, chatID, "Название слишком короткое. Пришли корректное название или /cancel.")
		s.sendMainMenu(ctx, chatID)
		return
	}
	if len([]rune(name)) > 64 {
		_ = s.tg.SendMessage(ctx, chatID, "Название слишком длинное (макс 64). Пришли короче или /cancel.")
		s.sendMainMenu(ctx, chatID)
		return
	}

	p, err := parsePayload(st.Payload)
	if err != nil {
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка payload. Напиши /addshop чтобы начать заново.")
		s.sendMainMenu(ctx, chatID)
		return
	}

	wbToken := strings.TrimSpace(p[PayloadWBToken])
	if wbToken == "" {
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Токен потерялся. Напиши /addshop чтобы начать заново.")
		s.sendMainMenu(ctx, chatID)
		return
	}

	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil {
		s.log.Error("ensure user failed", "chat_id", chatID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Не получилось создать пользователя. Попробуй ещё раз /addshop.")
		s.sendMainMenu(ctx, chatID)
		return
	}

	shop := &domain.Shop{
		UserID:      u.ID,
		Marketplace: "wildberries",
		ExternalID:  "",
		Name:        name,
		APIToken:    wbToken,
		IsActive:    true,
	}

	if err := s.shopsRepo.Create(ctx, shop); err != nil {
		s.log.Error("create shop failed", "chat_id", chatID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Не получилось создать магазин (ошибка БД). Попробуй ещё раз /addshop.")
		return
	}

	// 2) create telegram endpoint (owner)
	label := "owner"
	chatStr := fmt.Sprintf("%d", chatID)

	exists, err := s.endpointsRepo.Exists(ctx, shop.ID, "telegram", chatStr)
	if err != nil {
		s.log.Error("endpoint exists check failed", "shop_id", shop.ID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Магазин создан, но не получилось привязать чат. Напиши мне 🙂")
		return
	}
	if !exists {
		ep := &domain.NotificationEndpoint{
			ShopID:  shop.ID,
			Channel: "telegram",
			Address: chatStr,
			Enabled: true,
			Label:   &label,
		}
		if err := s.endpointsRepo.Create(ctx, ep); err != nil {
			s.log.Error("create endpoint failed", "shop_id", shop.ID, "err", err)
			_ = s.tg.SendMessage(ctx, chatID, "Магазин создан, но не получилось привязать чат. Напиши мне 🙂")
			return
		}
	}

	// 3) clear state
	_ = s.statesRepo.Clear(ctx, chatID)

	_ = s.tg.SendMessageWithKeyboard(ctx, chatID,
		"Готово ✅ Магазин подключен.\n\n"+
			"В ближайший час придёт первый sync/алерты (если они сработают).",
		removeKeyboard(),
	)
}

func parsePayload(b []byte) (payloadMap, error) {
	if len(b) == 0 {
		return payloadMap{}, nil
	}
	var p payloadMap
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	if p == nil {
		p = payloadMap{}
	}
	return p, nil
}

func mustJSON(p payloadMap) []byte {
	b, _ := json.Marshal(p)
	return b
}

func mustJSONAny(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
