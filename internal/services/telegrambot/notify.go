package telegrambot

import (
	"context"
	"fmt"
	"strconv"
)

func (s *Service) handleNotifyToggle(
	ctx context.Context,
	chatID int64,
	action string,
	shopIDStr string,
) {
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный shop_id.")
		return
	}

	enabled := false
	switch action {
	case "on":
		enabled = true
	case "off":
		enabled = false
	default:
		_ = s.tg.SendMessage(ctx, chatID, "Используй: /notify on <id> или /notify off <id>")
		return
	}

	addr := strconv.FormatInt(chatID, 10)

	if err := s.endpointsRepo.SetEnabledByShopAndAddress(
		ctx, shopID, "telegram", addr, enabled,
	); err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка обновления уведомлений.")
		return
	}

	status := "выключены"
	if enabled {
		status = "включены"
	}

	_ = s.tg.SendMessageWithKeyboard(ctx, chatID,
		fmt.Sprintf("Уведомления %s для магазина %d.", status, shopID),
		removeKeyboard(),
	)
}
