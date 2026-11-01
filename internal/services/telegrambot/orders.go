package telegrambot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const ordersCooldown = 5 * time.Minute

// handleOrders — команда /orders: показывает текущие заказы на сборке по всем магазинам чата.
// Данные берутся напрямую из WB API (GET /api/v3/orders/new) — всегда актуальные.
func (s *Service) handleOrders(ctx context.Context, chatID int64) {
	if s.wbClient == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Функция временно недоступна.")
		return
	}

	// rate limit
	now := time.Now()
	if v, ok := s.ordersRateLimit.Load(chatID); ok {
		last := v.(time.Time)
		if now.Sub(last) < ordersCooldown {
			remaining := ordersCooldown - now.Sub(last)
			mins := int(remaining.Minutes()) + 1
			_ = s.tg.SendMessage(ctx, chatID,
				fmt.Sprintf("Запрос недавно выполнялся. Повтори через %d мин.", mins))
			return
		}
	}
	s.ordersRateLimit.Store(chatID, now)

	shops, err := s.shopsRepo.ListByTelegramChat(ctx, chatID)
	if err != nil {
		s.log.Error("orders: list shops by telegram chat failed", "chat_id", chatID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка загрузки магазинов. Попробуй позже.")
		return
	}
	if len(shops) == 0 {
		s.log.Warn("orders: no shops found for chat", "chat_id", chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Не найдено привязанных магазинов.")
		return
	}

	var lines []string
	totalCount := 0

	for _, shop := range shops {
		fbsResp, err := s.wbClient.GetNewOrders(ctx, shop.APIToken)
		if err != nil {
			s.log.Warn("orders: FBS WB API failed", "shop_id", shop.ID, "err", err)
			lines = append(lines, fmt.Sprintf("❌ %s: ошибка запроса к WB (FBS)", shop.Name))
			continue
		}
		dbwResp, err := s.wbClient.GetNewDBWOrders(ctx, shop.APIToken)
		if err != nil {
			s.log.Warn("orders: DBW WB API failed", "shop_id", shop.ID, "err", err)
			// DBW failure is non-fatal — show FBS results anyway
			dbwResp = nil
		}

		fbsCount := len(fbsResp.Orders)
		dbwCount := 0
		if dbwResp != nil {
			dbwCount = len(dbwResp.Orders)
		}
		if fbsCount+dbwCount == 0 {
			continue
		}

		// Названия товаров из нашей БД (для подстановки)
		titleByNmID := make(map[int64]string)
		if s.productsRepo != nil {
			products, err := s.productsRepo.ListByShop(ctx, shop.ID, 0, 0)
			if err == nil {
				for _, p := range products {
					if p == nil || p.ExternalID == "" {
						continue
					}
					nmID, parseErr := strconv.ParseInt(p.ExternalID, 10, 64)
					if parseErr == nil {
						titleByNmID[nmID] = p.Title
					}
				}
			}
		}

		lines = append(lines, fmt.Sprintf("🏪 *%s* (%d шт.):", shop.Name, fbsCount+dbwCount))
		for _, o := range fbsResp.Orders {
			title := titleByNmID[o.NmID]
			if title == "" {
				title = o.Article
			}
			if title == "" {
				title = fmt.Sprintf("nm %d", o.NmID)
			}
			lines = append(lines, fmt.Sprintf("• %s — заказ #%d", title, o.ID))
		}
		if dbwResp != nil {
			for _, o := range dbwResp.Orders {
				title := titleByNmID[o.NmID]
				if title == "" {
					title = o.Article
				}
				if title == "" {
					title = fmt.Sprintf("nm %d", o.NmID)
				}
				lines = append(lines, fmt.Sprintf("• %s — заказ #%d (WB доставка)", title, o.ID))
			}
		}
		totalCount += fbsCount + dbwCount
	}

	if totalCount == 0 {
		_ = s.tg.SendMessage(ctx, chatID, "Нет активных заказов на сборке.")
		return
	}

	header := fmt.Sprintf("📦 Заказы на сборке (%d):\n\n", totalCount)
	_ = s.tg.SendMessage(ctx, chatID, strings.TrimRight(header+strings.Join(lines, "\n"), "\n"))
}
