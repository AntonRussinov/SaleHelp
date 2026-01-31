package alert

import (
	"context"
	"salehelp/domain"
	notifications "salehelp/internal/services/notification"
	"strconv"
	"time"
)

// runShopSaleEvents отправляет уведомления о новых выкупах/возвратах в Telegram и помечает события как уведомлённые.
func (e *AlertEngine) runShopSaleEvents(ctx context.Context, shopID int, products []*domain.Product) error {
	events, err := e.saleEventsRepo.ListUnnotifiedByShop(ctx, shopID, 50)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	if len(products) == 0 {
		products, err = e.productsRepo.ListByShop(ctx, shopID, 0, 0)
		if err != nil {
			return err
		}
	}
	nmIDSet := make(map[int64]struct{}, len(events))
	for _, ev := range events {
		nmIDSet[ev.NmID] = struct{}{}
	}
	titleByNmID := make(map[int64]string)
	for _, p := range products {
		if p == nil || p.Title == "" || p.ExternalID == "" {
			continue
		}
		nmID, err := strconv.ParseInt(p.ExternalID, 10, 64)
		if err != nil {
			continue
		}
		if _, ok := nmIDSet[nmID]; ok {
			titleByNmID[nmID] = p.Title
		}
	}
	endpoints, err := e.endpointsRepo.ListByShopID(ctx, shopID, true)
	if err != nil {
		return err
	}
	text := notifications.FormatSaleEventsMessageRU(events, titleByNmID)
	now := time.Now()
	sent := false
	for _, ep := range endpoints {
		if ep.Channel != "telegram" || !ep.Enabled {
			continue
		}
		chatID, err := strconv.ParseInt(ep.Address, 10, 64)
		if err != nil {
			e.log.Warn("sale events: invalid telegram chat_id", "endpoint_id", ep.ID, "address", ep.Address, "err", err)
			continue
		}
		if err := e.tg.SendMessage(ctx, chatID, text); err != nil {
			e.log.Error("sale events: send failed", "endpoint_id", ep.ID, "err", err)
			continue
		}
		sent = true
	}
	if sent {
		ids := make([]int, 0, len(events))
		for _, ev := range events {
			ids = append(ids, ev.ID)
		}
		if err := e.saleEventsRepo.MarkNotified(ctx, ids, now); err != nil {
			return err
		}
	}
	return nil
}
