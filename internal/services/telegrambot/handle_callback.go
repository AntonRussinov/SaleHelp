package telegrambot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"salehelp/domain"
	"salehelp/internal/handler/telegram/dto"
)

const (
	callbackPrefixABCShop        = "abc_shop_"
	callbackPrefixRestorePrices  = "restore_prices:"
	callbackPrefixConfirmRestore = "confirm_restore:"
)

func (s *Service) handleCallbackQuery(ctx context.Context, q *dto.CallbackQuery) {
	if q == nil {
		return
	}

	chatID := int64(0)
	if q.Message != nil {
		chatID = q.Message.Chat.ID
	}
	if chatID == 0 {
		return
	}

	_ = s.tg.AnswerCallbackQuery(ctx, q.ID)

	data := strings.TrimSpace(q.Data)

	if strings.HasPrefix(data, callbackPrefixABCShop) {
		shopIDStr := strings.TrimPrefix(data, callbackPrefixABCShop)
		shopID, err := strconv.Atoi(shopIDStr)
		if err != nil {
			_ = s.tg.SendMessage(ctx, chatID, "Неверный выбор магазина.")
			return
		}
		s.handleABCNow(ctx, chatID, []string{"/abc_now", strconv.Itoa(shopID)})
		return
	}

	if strings.HasPrefix(data, callbackPrefixRestorePrices) {
		msgID := 0
		if q.Message != nil {
			msgID = q.Message.MessageID
		}
		s.handleRestorePrices(ctx, chatID, msgID, strings.TrimPrefix(data, callbackPrefixRestorePrices))
		return
	}

	if strings.HasPrefix(data, callbackPrefixConfirmRestore) {
		s.handleConfirmRestore(ctx, chatID, strings.TrimPrefix(data, callbackPrefixConfirmRestore))
		return
	}

	if strings.HasPrefix(data, callbackPrefixAlertSettingsShop) {
		shopIDStr := strings.TrimPrefix(data, callbackPrefixAlertSettingsShop)
		shopID, err := strconv.Atoi(shopIDStr)
		if err != nil {
			_ = s.tg.SendMessage(ctx, chatID, "Неверный выбор магазина.")
			return
		}
		s.handleAlertSettingsShop(ctx, chatID, shopID)
		return
	}

	if strings.HasPrefix(data, callbackPrefixAlertToggle) {
		payload := strings.TrimPrefix(data, callbackPrefixAlertToggle)
		parts := strings.SplitN(payload, ":", 2)
		if len(parts) != 2 {
			_ = s.tg.SendMessage(ctx, chatID, "Неверный запрос.")
			return
		}
		shopID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			_ = s.tg.SendMessage(ctx, chatID, "Неверный shop_id.")
			return
		}
		key := strings.TrimSpace(parts[1])
		s.handleAlertToggle(ctx, chatID, shopID, key)
		return
	}

	if strings.HasPrefix(data, callbackPrefixAlertAll) {
		payload := strings.TrimPrefix(data, callbackPrefixAlertAll)
		parts := strings.SplitN(payload, ":", 2)
		if len(parts) != 2 {
			_ = s.tg.SendMessage(ctx, chatID, "Неверный запрос.")
			return
		}
		shopID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			_ = s.tg.SendMessage(ctx, chatID, "Неверный shop_id.")
			return
		}
		action := strings.TrimSpace(parts[1])
		s.handleAlertAll(ctx, chatID, shopID, action)
		return
	}
}

func (s *Service) handleRestorePrices(ctx context.Context, chatID int64, messageID int, payload string) {
	batchID, err := strconv.Atoi(strings.TrimSpace(payload))
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный идентификатор запроса.")
		return
	}
	if s.restoreBatchRepo != nil {
		batch, _ := s.restoreBatchRepo.GetByID(ctx, batchID)
		if batch != nil {
			switch batch.RestoreStatus {
			case domain.RestoreStatusRestoring:
				_ = s.tg.SendMessage(ctx, chatID, "⏳ Цены уже восстанавливаются, подождите...")
				return
			case domain.RestoreStatusRestored:
				t := batch.RestoredAt
				if t == nil {
					_ = s.tg.SendMessage(ctx, chatID, "ℹ️ Цены уже были восстановлены.")
				} else {
					msk := t.In(time.FixedZone("MSK", 3*3600))
					_ = s.tg.SendMessage(ctx, chatID, "ℹ️ Цены уже были восстановлены в "+msk.Format("15:04")+" МСК")
				}
				return
			case domain.RestoreStatusFailed:
				// allow retry — fall through to confirm dialog
			}
		}
	}

	// Edit the original alert message keyboard in place so repeated clicks don't spawn new dialogs.
	confirmBase := fmt.Sprintf("%s%d:%d:", callbackPrefixConfirmRestore, batchID, messageID)
	keyboard := map[string]any{
		"inline_keyboard": [][]map[string]string{
			{
				{"text": "✅ Да, восстановить", "callback_data": confirmBase + "yes"},
				{"text": "❌ Нет, оставить", "callback_data": confirmBase + "no"},
			},
		},
	}
	_ = s.tg.EditMessageReplyMarkup(ctx, chatID, messageID, keyboard)
}

func (s *Service) handleConfirmRestore(ctx context.Context, chatID int64, payload string) {
	// payload format: "<batchID>:<messageID>:<yes|no>"
	parts := strings.SplitN(payload, ":", 3)
	if len(parts) != 3 {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный запрос.")
		return
	}
	batchID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный идентификатор запроса.")
		return
	}
	messageID, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	choice := strings.TrimSpace(strings.ToLower(parts[2]))

	restoreOriginalKeyboard := func() {
		kb := map[string]any{
			"inline_keyboard": [][]map[string]string{
				{{"text": "🔄 Вернуть цены назад", "callback_data": callbackPrefixRestorePrices + strconv.Itoa(batchID)}},
			},
		}
		_ = s.tg.EditMessageReplyMarkup(ctx, chatID, messageID, kb)
	}

	if choice == "no" {
		restoreOriginalKeyboard()
		return
	}
	if choice != "yes" {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный запрос.")
		return
	}
	if s.priceRestore == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Сервис восстановления цен недоступен.")
		restoreOriginalKeyboard()
		return
	}
	result := s.priceRestore.Restore(ctx, batchID)
	if result.AlreadyRestoring {
		_ = s.tg.SendMessage(ctx, chatID, "⏳ Цены уже восстанавливаются, подождите...")
		return
	}
	if result.AlreadyRestored && result.RestoredAt != nil {
		t := result.RestoredAt.In(time.FixedZone("MSK", 3*3600))
		_ = s.tg.EditMessageReplyMarkup(ctx, chatID, messageID, map[string]any{"inline_keyboard": [][]map[string]string{}})
		_ = s.tg.SendMessage(ctx, chatID, "ℹ️ Цены уже были восстановлены в "+t.Format("15:04")+" МСК")
		return
	}
	if result.Error != nil {
		restoreOriginalKeyboard()
		_ = s.tg.SendMessage(ctx, chatID,
			"❌ Не удалось восстановить цены\n\nОшибка: "+result.Error.Error()+"\n\nПопробуйте вручную или повторите через несколько минут.")
		return
	}

	// Remove keyboard from the original alert message — restore is done.
	_ = s.tg.EditMessageReplyMarkup(ctx, chatID, messageID, map[string]any{"inline_keyboard": [][]map[string]string{}})

	msk := time.FixedZone("MSK", 3*3600)
	now := time.Now().In(msk)
	successMsg := fmt.Sprintf("✅ Цены восстановлены!\nВосстановлено товаров: %d\nВремя: %s МСК", result.RestoredCount, now.Format("15:04"))

	// broadcast to all enabled Telegram endpoints for the shop
	if result.ShopID > 0 && s.endpointsRepo != nil {
		endpoints, err := s.endpointsRepo.ListByShopID(ctx, result.ShopID, true)
		if err == nil {
			notified := map[int64]struct{}{}
			for _, ep := range endpoints {
				if ep == nil || ep.Channel != "telegram" {
					continue
				}
				cid, errParse := strconv.ParseInt(ep.Address, 10, 64)
				if errParse != nil {
					continue
				}
				if _, already := notified[cid]; already {
					continue
				}
				notified[cid] = struct{}{}
				_ = s.tg.SendMessage(ctx, cid, successMsg)
			}
			return
		}
	}
	// fallback: reply only to the requesting chat
	_ = s.tg.SendMessage(ctx, chatID, successMsg)
}
