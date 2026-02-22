package telegrambot

import (
	"context"
	"strconv"
	"strings"

	"salehelp/internal/handler/telegram/dto"
)

const callbackPrefixABCShop = "abc_shop_"

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
}
