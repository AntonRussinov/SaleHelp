package telegrambot

import (
	"context"
	"fmt"
	"salehelp/domain"
	"strconv"
	"strings"
)

func (s *Service) handleMyShops(ctx context.Context, chatID int64) error {
	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка пользователя.")
		return nil
	}

	// 1) магазины, где user — owner
	ownerShops, err := s.shopsRepo.ListByUser(ctx, u.ID, 0, 0)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Не удалось получить магазины.")
		return nil
	}

	// 2) магазины, где user подключён как telegram-участник
	memberShops, err := s.shopsRepo.ListByTelegramChat(ctx, chatID)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Не удалось получить магазины.")
		return nil
	}

	// 3) объединяем без дублей
	shopsByID := make(map[int]*domain.Shop)

	for _, sh := range ownerShops {
		shopsByID[sh.ID] = sh
	}
	for _, sh := range memberShops {
		if _, ok := shopsByID[sh.ID]; !ok {
			shopsByID[sh.ID] = sh
		}
	}

	if len(shopsByID) == 0 {
		_ = s.tg.SendMessage(ctx, chatID, "У тебя пока нет магазинов.")
		s.sendMainMenu(ctx, chatID)
		return nil
	}

	addr := strconv.FormatInt(chatID, 10)
	var out strings.Builder

	out.WriteString("📦 Твои магазины:\n\n")

	for _, sh := range shopsByID {
		endpoint, _ := s.endpointsRepo.GetByShopAndAddress(
			ctx, sh.ID, "telegram", addr,
		)

		notifyStatus := "OFF"
		if endpoint != nil && endpoint.Enabled {
			notifyStatus = "ON"
		}

		shopStatus := "OFF"
		if sh.IsActive {
			shopStatus = "ON"
		}

		role := "member"
		if sh.UserID == u.ID {
			role = "owner"
		}

		out.WriteString(fmt.Sprintf(
			"• %s (id=%d)\n  Роль: %s\n  Магазин: %s\n  Уведомления тебе: %s\n",
			sh.Name, sh.ID, role, shopStatus, notifyStatus,
		))

		out.WriteString("  Команды:\n")

		// owner может выключать магазин глобально
		if role == "owner" {
			out.WriteString(fmt.Sprintf(
				"  /shop off %d\n  /shop on %d\n",
				sh.ID, sh.ID,
			))
			out.WriteString(fmt.Sprintf(
				"  /chats %d\n",
				sh.ID,
			))
		}

		// любой пользователь управляет своими уведомлениями
		out.WriteString(fmt.Sprintf(
			"  /notify off %d\n  /notify on %d\n",
			sh.ID, sh.ID,
		))

		out.WriteString("\n")
	}

	_ = s.tg.SendMessage(ctx, chatID, out.String())
	s.sendMainMenu(ctx, chatID)
	return nil
}
