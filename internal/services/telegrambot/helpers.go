package telegrambot

import (
	"context"
	"encoding/json"
	"fmt"
	"salehelp/domain"
)

func encodeInts(xs []int) string {
	b, _ := json.Marshal(xs)
	return string(b)
}

func decodeInts(s string) ([]int, error) {
	var xs []int
	if s == "" {
		return xs, nil
	}
	err := json.Unmarshal([]byte(s), &xs)
	return xs, err
}

func extractShopIDs(shops []*domain.Shop) []int {
	ids := make([]int, 0, len(shops))
	for _, sh := range shops {
		ids = append(ids, sh.ID)
	}
	return ids
}

func mainMenuKeyboard() map[string]any {
	return map[string]any{
		"keyboard": [][]map[string]string{
			{
				{"text": "/addshop"},
				{"text": "/myshops"},
			},
			{
				{"text": "/invite"},
			},
			{
				{"text": "/cancel"},
			},
		},
		"resize_keyboard": true,
	}
}

func (s *Service) sendMainMenu(ctx context.Context, chatID int64) {
	_ = s.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		"Главное меню:",
		mainMenuKeyboard(),
	)
}

func (s *Service) canInviteMoreTelegramChats(ctx context.Context, shopID int) (bool, string) {
	allowed, err := s.subEntitlements.AllowedTelegramChats(ctx, shopID)
	if err != nil {
		return false, "Не удалось проверить подписку. Попробуй позже."
	}

	// Политика v0: без подписки нельзя добавлять чаты/инвайтить
	if allowed <= 0 {
		return false, "У магазина нет активной подписки. Приглашения доступны только при активной подписке."
	}

	current, err := s.endpointsRepo.CountEnabledByShopAndChannel(ctx, shopID, "telegram")
	if err != nil {
		return false, "Не удалось проверить текущие подключенные чаты. Попробуй позже."
	}

	if current >= allowed {
		return false, fmt.Sprintf(
			"Лимит Telegram-чатов по подписке исчерпан: %d/%d.\n\nПродли подписку или подключи доп. чаты.",
			current, allowed,
		)
	}

	return true, ""
}
