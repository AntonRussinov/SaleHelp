package telegrambot

import (
	"context"
	"fmt"
	"salehelp/domain"
	"strconv"
	"strings"
	"time"
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

	// last_sync_at для всех магазинов пачкой
	shopIDs := make([]int, 0, len(shopsByID))
	for id := range shopsByID {
		shopIDs = append(shopIDs, id)
	}
	lastSyncMap, _ := s.shopsRepo.GetLastSyncAtMap(ctx, shopIDs)

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

		// --- NEW: subscription + addons + last sync ---
		// 1) подписка
		subLine := "нет"
		level, sub, subErr := s.subEntitlements.GetShopAccessLevel(ctx, sh.ID)
		if subErr == nil && sub != nil {
			// показываем valid_to и level
			subLine = fmt.Sprintf("%s, до %s", level, sub.ValidTo.Format("02.01.2006"))
			if level == "blocked" {
				subLine = fmt.Sprintf("%s (sync остановлен)", subLine)
			} else if level == "grace" {
				graceUntil := sub.ValidTo.Add(72 * time.Hour)
				subLine = fmt.Sprintf("истекла, grace до %s", graceUntil.Format("02.01.2006"))
			}
		}

		// 2) лимит чатов (1 базовый + допы)
		allowedChats, _ := s.subEntitlements.AllowedTelegramChats(ctx, sh.ID)
		addonChats := allowedChats - 1
		if addonChats < 0 {
			addonChats = 0
		}

		// 3) last sync
		lastSyncText := "нет данных"
		if t := lastSyncMap[sh.ID]; t != nil {
			lastSyncText = t.Format("02.01 15:04")
		}
		// --- /NEW ---

		out.WriteString(fmt.Sprintf(
			"• %s (id=%d)\n"+
				"  Роль: %s\n"+
				"  Магазин: %s\n"+
				"  Уведомления тебе: %s\n"+
				"  Подписка: %s\n"+
				"  Доп. чаты: %d (лимит чатов: %d)\n"+
				"  Последний синк: %s\n",
			sh.Name, sh.ID,
			role,
			shopStatus,
			notifyStatus,
			subLine,
			addonChats, allowedChats,
			lastSyncText,
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
