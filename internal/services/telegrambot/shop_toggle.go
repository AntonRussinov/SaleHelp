package telegrambot

import (
	"context"
	"fmt"
	"strconv"
)

func (s *Service) handleShopToggle(
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

	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка пользователя.")
		return
	}

	shop, err := s.shopsRepo.GetByID(ctx, shopID)
	if err != nil || shop == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Магазин не найден.")
		return
	}

	if shop.UserID != u.ID {
		_ = s.tg.SendMessage(ctx, chatID, "У тебя нет прав на этот магазин.")
		return
	}

	var active bool
	switch action {
	case "on":
		active = true
	case "off":
		active = false
	default:
		_ = s.tg.SendMessage(ctx, chatID, "Используй: /shop on <id> или /shop off <id>")
		return
	}

	if err := s.shopsRepo.SetActive(ctx, shopID, active); err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка обновления магазина.")
		return
	}

	status := "включён"
	if !active {
		status = "выключен"
	}

	_ = s.tg.SendMessageWithKeyboard(ctx, chatID,
		fmt.Sprintf("Магазин \"%s\" теперь %s.", shop.Name, status),
		removeKeyboard(),
	)
}
