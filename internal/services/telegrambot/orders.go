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
	if err != nil || len(shops) == 0 {
		_ = s.tg.SendMessage(ctx, chatID, "Не найдено привязанных магазинов.")
		return
	}

	var lines []string
	totalCount := 0

	for _, shop := range shops {
		resp, err := s.wbClient.GetNewOrders(ctx, shop.APIToken)
		if err != nil {
			s.log.Warn("orders: WB API failed", "shop_id", shop.ID, "err", err)
			lines = append(lines, fmt.Sprintf("❌ %s: ошибка запроса к WB", shop.Name))
			continue
		}
		if len(resp.Orders) == 0 {
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

		lines = append(lines, fmt.Sprintf("🏪 *%s* (%d шт.):", shop.Name, len(resp.Orders)))
		for _, o := range resp.Orders {
			title := titleByNmID[o.NmID]
			if title == "" {
				title = o.Article
			}
			if title == "" {
				title = fmt.Sprintf("nm %d", o.NmID)
			}
			lines = append(lines, fmt.Sprintf("• %s — заказ #%d", title, o.ID))
		}
		totalCount += len(resp.Orders)
	}

	if totalCount == 0 {
		_ = s.tg.SendMessage(ctx, chatID, "Нет активных заказов на сборке.")
		return
	}

	header := fmt.Sprintf("📦 Заказы на сборке (%d):\n\n", totalCount)
	_ = s.tg.SendMessage(ctx, chatID, strings.TrimRight(header+strings.Join(lines, "\n"), "\n"))
}
