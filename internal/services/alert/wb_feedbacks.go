package alert

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"salehelp/domain"
)

// runShopWBFeedbacks отправляет уведомления о новых отзывах в Telegram и помечает как уведомлённые.
func (e *AlertEngine) runShopWBFeedbacks(ctx context.Context, shopID int) error {
	feedbacks, err := e.wbFeedbacksRepo.ListUnnotifiedByShop(ctx, shopID, 50)
	if err != nil {
		return err
	}
	if len(feedbacks) == 0 {
		return nil
	}
	endpoints, err := e.endpointsRepo.ListByShopID(ctx, shopID, true)
	if err != nil {
		return err
	}
	now := time.Now()
	sent := false
	for _, f := range feedbacks {
		text := formatWBFeedbackMessageRU(f)
		for _, ep := range endpoints {
			if ep.Channel != "telegram" || !ep.Enabled {
				continue
			}
			chatID, err := strconv.ParseInt(ep.Address, 10, 64)
			if err != nil {
				e.log.Warn("wb feedbacks: invalid telegram chat_id", "endpoint_id", ep.ID, "address", ep.Address, "err", err)
				continue
			}
			if err := e.tg.SendMessage(ctx, chatID, text); err != nil {
				e.log.Error("wb feedbacks: send failed", "endpoint_id", ep.ID, "feedback_id", f.ID, "err", err)
				continue
			}
			sent = true
		}
	}
	if sent {
		ids := make([]int, 0, len(feedbacks))
		for _, f := range feedbacks {
			ids = append(ids, f.ID)
		}
		if err := e.wbFeedbacksRepo.MarkNotified(ctx, ids, now); err != nil {
			return err
		}
	}
	return nil
}

// formatWBFeedbackMessageRU форматирует сообщение о новом отзыве для Telegram.
func formatWBFeedbackMessageRU(f *domain.WBFeedback) string {
	productName := f.ProductName
	if productName == "" {
		productName = fmt.Sprintf("nmId %d", f.NmID)
	}
	return fmt.Sprintf("⭐ Новый отзыв о товаре\n\nТовар: %s\nОтзыв: \"%s\"",
		productName, f.Text)
}
