package telegrambot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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
		s.handleRestorePrices(ctx, chatID, strings.TrimPrefix(data, callbackPrefixRestorePrices))
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

func (s *Service) handleRestorePrices(ctx context.Context, chatID int64, payload string) {
	batchID, err := strconv.Atoi(strings.TrimSpace(payload))
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный идентификатор запроса.")
		return
	}
	if s.restoreBatchRepo != nil {
		batch, _ := s.restoreBatchRepo.GetByID(ctx, batchID)
		if batch != nil && batch.RestoredAt != nil {
			t := batch.RestoredAt.In(time.FixedZone("MSK", 3*3600))
			_ = s.tg.SendMessage(ctx, chatID, "ℹ️ Цены уже были восстановлены в "+t.Format("15:04")+" МСК")
			return
		}
	}
	// Отправляем сообщение с подтверждением и кнопками
	text := "Восстановить цены на прежний уровень по этим товарам?"
	keyboard := map[string]any{
		"inline_keyboard": [][]map[string]string{
			{
				{"text": "Да, изменить", "callback_data": callbackPrefixConfirmRestore + strconv.Itoa(batchID) + ":yes"},
				{"text": "Нет, оставить как есть", "callback_data": callbackPrefixConfirmRestore + strconv.Itoa(batchID) + ":no"},
			},
		},
	}
	_ = s.tg.SendMessageWithInlineKeyboard(ctx, chatID, text, keyboard)
}

func (s *Service) handleConfirmRestore(ctx context.Context, chatID int64, payload string) {
	parts := strings.SplitN(payload, ":", 2)
	if len(parts) != 2 {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный запрос.")
		return
	}
	batchID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный идентификатор запроса.")
		return
	}
	choice := strings.TrimSpace(strings.ToLower(parts[1]))
	if choice == "no" {
		_ = s.tg.SendMessage(ctx, chatID, "Запрос отменён.")
		return
	}
	if choice != "yes" {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный запрос.")
		return
	}
	if s.priceRestore == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Сервис восстановления цен недоступен.")
		return
	}
	result := s.priceRestore.Restore(ctx, batchID)
	if result.AlreadyRestored && result.RestoredAt != nil {
		t := result.RestoredAt.In(time.FixedZone("MSK", 3*3600))
		_ = s.tg.SendMessage(ctx, chatID, "ℹ️ Цены уже были восстановлены в "+t.Format("15:04")+" МСК")
		return
	}
	if result.Error != nil {
		_ = s.tg.SendMessage(ctx, chatID,
			"❌ Не удалось восстановить цены\n\nОшибка: "+result.Error.Error()+"\n\nПопробуйте вручную или повторите через несколько минут.")
		return
	}
	msk := time.FixedZone("MSK", 3*3600)
	now := time.Now().In(msk)
	_ = s.tg.SendMessage(ctx, chatID,
		fmt.Sprintf("✅ Цены восстановлены!\nВосстановлено товаров: %d\nВремя: %s МСК", result.RestoredCount, now.Format("15:04")))
}
