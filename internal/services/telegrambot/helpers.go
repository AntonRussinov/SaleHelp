package telegrambot

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

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

// inlineKeyboardABCShops строит inline-клавиатуру "выбор магазина для ABC".
// callback_data: "abc_shop_<id>".
func inlineKeyboardABCShops(shopsByID map[int]*domain.Shop) map[string]any {
	// Упорядочиваем по ID для стабильного отображения
	ids := make([]int, 0, len(shopsByID))
	for id := range shopsByID {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	row := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		sh := shopsByID[id]
		name := sh.Name
		if name == "" {
			name = "Магазин " + strconv.Itoa(id)
		}
		row = append(row, map[string]string{
			"text":          name,
			"callback_data": "abc_shop_" + strconv.Itoa(id),
		})
	}
	return map[string]any{
		"inline_keyboard": [][]map[string]string{row},
	}
}

// removeKeyboard — убирает reply-клавиатуру. Команды доступны через меню бота (☰).
func removeKeyboard() map[string]any {
	return map[string]any{
		"remove_keyboard": true,
	}
}

func (s *Service) sendMainMenu(ctx context.Context, chatID int64) {
	_ = s.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		"Команды доступны в меню (☰ рядом с полем ввода).",
		removeKeyboard(),
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

func formatMinor(amount int64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}
	return fmt.Sprintf("%s%d.%02d", sign, amount/100, amount%100)
}
