package telegrambot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

func (s *Service) handleChats(ctx context.Context, chatID int64, shopIDStr string) {
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный shop_id.")
		return
	}

	u, _ := s.ensureUserByChatID(ctx, chatID)
	shop, _ := s.shopsRepo.GetByID(ctx, shopID)

	if shop == nil || shop.UserID != u.ID {
		_ = s.tg.SendMessage(ctx, chatID, "Нет доступа.")
		return
	}

	endpoints, err := s.endpointsRepo.ListByShopID(ctx, shopID, false)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка получения чатов.")
		return
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("📦 Чаты магазина \"%s\":\n\n", shop.Name))

	for i, e := range endpoints {
		status := "OFF"
		if e.Enabled {
			status = "ON"
		}
		out.WriteString(fmt.Sprintf(
			"%d) chat %s — %s\n",
			i+1, e.Address, status,
		))
	}

	out.WriteString("\nКоманды:\n/chat on <shop_id> <номер>\n/chat off <shop_id> <номер>")

	_ = s.tg.SendMessage(ctx, chatID, out.String())
}
