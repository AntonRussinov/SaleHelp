package telegrambot

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"salehelp/domain"
	"strconv"
	"strings"
	"time"
)

func (s *Service) startInvite(ctx context.Context, chatID int64) {
	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка пользователя.")
		return
	}

	shops, err := s.shopsRepo.ListByUser(ctx, u.ID, 0, 0)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Не удалось получить магазины.")
		return
	}

	if len(shops) == 0 {
		_ = s.tg.SendMessage(ctx, chatID, "У тебя нет магазинов для приглашения.")
		return
	}

	var text = "Выбери магазин для приглашения:\n\n"
	for i, sh := range shops {
		text += fmt.Sprintf("%d) %s\n", i+1, sh.Name)
	}

	type invitePayload struct {
		ShopIDs []int `json:"shop_ids"`
	}

	state := &domain.TelegramChatState{
		ChatID: chatID,
		State:  StateInviteChooseShop,
		Payload: mustJSONAny(invitePayload{
			ShopIDs: extractShopIDs(shops),
		}),
	}

	_ = s.statesRepo.Upsert(ctx, state, time.Now())
	_ = s.tg.SendMessage(ctx, chatID, text)
}

func (s *Service) handleInviteChooseShop(
	ctx context.Context,
	chatID int64,
	text string,
	st *domain.TelegramChatState,
) {
	idx, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil || idx <= 0 {
		_ = s.tg.SendMessage(ctx, chatID, "Введите номер магазина.")
		return
	}

	var data struct {
		ShopIDs []int `json:"shop_ids"`
	}

	if err := json.Unmarshal(st.Payload, &data); err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка данных. Напиши /invite ещё раз.")
		_ = s.statesRepo.Clear(ctx, chatID)
		return
	}

	if idx > len(data.ShopIDs) {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный номер.")
		return
	}

	shopID := data.ShopIDs[idx-1]

	existing, err := s.invitesRepo.GetActiveByShopID(ctx, shopID)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка при проверке приглашений.")
		return
	}

	var invite *domain.ShopInviteToken

	if existing != nil {
		invite = existing
	} else {
		token := generateInviteToken()

		invite = &domain.ShopInviteToken{
			Token:     token,
			ShopID:    shopID,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		if err := s.invitesRepo.Create(ctx, invite); err != nil {
			_ = s.tg.SendMessage(ctx, chatID, "Ошибка создания токена.")
			return
		}
	}

	_ = s.statesRepo.Clear(ctx, chatID)

	_ = s.tg.SendMessage(ctx, chatID,
		fmt.Sprintf(
			"Готово ✅\n\nТокен для подключения:\n<code>%s</code>\n\n"+
				"Он действует до %s\n\n"+
				"Передай его и пусть напишут:\n/join %s",
			invite.Token,
			invite.ExpiresAt.Format("02.01 15:04"),
			invite.Token,
		),
	)

}

func generateInviteToken() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("INV-%x", b)
}
