package alert

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"salehelp/domain"
)

// runShopWBQuestions отправляет уведомления о новых вопросах в Telegram и помечает как уведомлённые.
func (e *AlertEngine) runShopWBQuestions(ctx context.Context, shopID int) error {
	questions, err := e.wbQuestionsRepo.ListUnnotifiedByShop(ctx, shopID, 50)
	if err != nil {
		return err
	}
	if len(questions) == 0 {
		return nil
	}
	endpoints, err := e.endpointsRepo.ListByShopID(ctx, shopID, true)
	if err != nil {
		return err
	}
	now := time.Now()
	sent := false
	for _, q := range questions {
		text := formatWBQuestionMessageRU(q)
		for _, ep := range endpoints {
			if ep.Channel != "telegram" || !ep.Enabled {
				continue
			}
			chatID, err := strconv.ParseInt(ep.Address, 10, 64)
			if err != nil {
				e.log.Warn("wb questions: invalid telegram chat_id", "endpoint_id", ep.ID, "address", ep.Address, "err", err)
				continue
			}
			if err := e.tg.SendMessage(ctx, chatID, text); err != nil {
				e.log.Error("wb questions: send failed", "endpoint_id", ep.ID, "question_id", q.ID, "err", err)
				continue
			}
			sent = true
		}
	}
	if sent {
		ids := make([]int, 0, len(questions))
		for _, q := range questions {
			ids = append(ids, q.ID)
		}
		if err := e.wbQuestionsRepo.MarkNotified(ctx, ids, now); err != nil {
			return err
		}
	}
	return nil
}

// formatWBQuestionMessageRU форматирует сообщение о новом вопросе для Telegram.
func formatWBQuestionMessageRU(q *domain.WBQuestion) string {
	productName := q.ProductName
	if productName == "" {
		productName = fmt.Sprintf("nmId %d", q.NmID)
	}
	return fmt.Sprintf("❓ Новый вопрос о товаре\n\nТовар: %s\nВопрос: \"%s\"",
		productName, q.Text)
}
