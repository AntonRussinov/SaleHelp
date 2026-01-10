package telegrambot

import (
	"context"
	"encoding/json"
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
