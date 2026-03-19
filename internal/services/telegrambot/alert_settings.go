package telegrambot

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"salehelp/domain"
	"salehelp/internal/notificationevents"
)

const (
	callbackPrefixAlertSettingsShop = "alert_settings_shop:"
	callbackPrefixAlertToggle       = "alert_toggle:"
	callbackPrefixAlertAll          = "alert_all:"
)

func (s *Service) handleAlertSettingsEntry(ctx context.Context, chatID int64) error {
	if s.notifSettings == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Настройки уведомлений временно недоступны.")
		return nil
	}

	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil || u == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка пользователя.")
		return nil
	}

	shops, err := s.shopsRepo.ListByUser(ctx, u.ID, 0, 0)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Не удалось получить магазины.")
		return nil
	}
	if len(shops) == 0 {
		_ = s.tg.SendMessage(ctx, chatID, "У тебя нет магазинов (owner).")
		return nil
	}

	shopsByID := make(map[int]*domain.Shop, len(shops))
	for _, sh := range shops {
		shopsByID[sh.ID] = sh
	}

	kb := inlineKeyboardShopsWithPrefix(shopsByID, callbackPrefixAlertSettingsShop)
	_ = s.tg.SendMessageWithInlineKeyboard(ctx, chatID, "Выбери магазин для настройки уведомлений:", kb)
	return nil
}

func (s *Service) handleAlertSettingsShop(ctx context.Context, chatID int64, shopID int) {
	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil || u == nil {
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

	s.sendAlertSettingsMenu(ctx, chatID, shop)
}

func (s *Service) handleAlertToggle(ctx context.Context, chatID int64, shopID int, key string) {
	if s.notifSettings == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Настройки уведомлений временно недоступны.")
		return
	}
	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil || u == nil {
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
	if _, ok := notificationevents.Find(key); !ok {
		_ = s.tg.SendMessage(ctx, chatID, "Неизвестный тип уведомлений.")
		return
	}

	cur, err := s.notifSettings.EffectiveEnabled(ctx, shopID, key)
	if err != nil {
		cur = true
	}
	_ = s.notifSettings.Set(ctx, shopID, key, !cur)

	s.sendAlertSettingsMenu(ctx, chatID, shop)
}

func (s *Service) handleAlertAll(ctx context.Context, chatID int64, shopID int, action string) {
	if s.notifSettings == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Настройки уведомлений временно недоступны.")
		return
	}
	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil || u == nil {
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

	enabled := true
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "on":
		enabled = true
	case "off":
		enabled = false
	default:
		_ = s.tg.SendMessage(ctx, chatID, "Неверная команда.")
		return
	}

	for _, def := range notificationevents.All() {
		_ = s.notifSettings.Set(ctx, shopID, def.Key, enabled)
	}

	s.sendAlertSettingsMenu(ctx, chatID, shop)
}

func (s *Service) sendAlertSettingsMenu(ctx context.Context, chatID int64, shop *domain.Shop) {
	if s.notifSettings == nil || shop == nil {
		return
	}

	m, err := s.notifSettings.EffectiveMap(ctx, shop.ID)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Не удалось получить настройки.")
		return
	}

	defs := notificationevents.All()
	// stable order by category then title
	sort.SliceStable(defs, func(i, j int) bool {
		if defs[i].Category != defs[j].Category {
			return defs[i].Category < defs[j].Category
		}
		return defs[i].Title < defs[j].Title
	})

	lines := make([]string, 0, len(defs)+2)
	lines = append(lines, fmt.Sprintf("⚙️ Уведомления магазина \"%s\" (id=%d):", shop.Name, shop.ID), "")
	for _, d := range defs {
		enabled := d.DefaultEnabled
		if v, ok := m[d.Key]; ok {
			enabled = v
		}
		state := "ON"
		if !enabled {
			state = "OFF"
		}
		lines = append(lines, fmt.Sprintf("• %s — %s", d.Title, state))
	}

	kb := inlineKeyboardAlertSettings(shop.ID, defs, m)
	_ = s.tg.SendMessageWithInlineKeyboard(ctx, chatID, strings.Join(lines, "\n"), kb)
}

func inlineKeyboardAlertSettings(shopID int, defs []notificationevents.Event, enabledMap map[string]bool) map[string]any {
	rows := make([][]map[string]string, 0, len(defs)+2)

	for _, d := range defs {
		enabled := d.DefaultEnabled
		if v, ok := enabledMap[d.Key]; ok {
			enabled = v
		}
		prefix := "✅ "
		if !enabled {
			prefix = "❌ "
		}
		rows = append(rows, []map[string]string{
			{
				"text":          prefix + d.Title,
				"callback_data": callbackPrefixAlertToggle + strconv.Itoa(shopID) + ":" + d.Key,
			},
		})
	}

	rows = append(rows, []map[string]string{
		{"text": "Включить все", "callback_data": callbackPrefixAlertAll + strconv.Itoa(shopID) + ":on"},
		{"text": "Выключить все", "callback_data": callbackPrefixAlertAll + strconv.Itoa(shopID) + ":off"},
	})

	return map[string]any{"inline_keyboard": rows}
}

func inlineKeyboardShopsWithPrefix(shopsByID map[int]*domain.Shop, prefix string) map[string]any {
	ids := make([]int, 0, len(shopsByID))
	for id := range shopsByID {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	row := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		sh := shopsByID[id]
		name := sh.Name
		if strings.TrimSpace(name) == "" {
			name = "Магазин " + strconv.Itoa(id)
		}
		row = append(row, map[string]string{
			"text":          name,
			"callback_data": prefix + strconv.Itoa(id),
		})
	}
	return map[string]any{
		"inline_keyboard": [][]map[string]string{row},
	}
}
