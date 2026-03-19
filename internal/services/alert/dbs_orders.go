package alert

import (
	"context"
	"salehelp/domain"
	notifications "salehelp/internal/services/notification"
	"strconv"
	"time"
)

// runShopDBSOrders отправляет уведомления о новых DBS-заказах в Telegram и помечает заказы как уведомлённые.
func (e *AlertEngine) runShopDBSOrders(ctx context.Context, shopID int, products []*domain.Product) error {
	orders, err := e.dbsOrdersRepo.ListUnnotifiedByShop(ctx, shopID, 50)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		return nil
	}

	if !e.notifEnabled(ctx, shopID, domain.NotificationEventDBSOrders) {
		now := time.Now()
		ids := make([]int, 0, len(orders))
		for _, o := range orders {
			ids = append(ids, o.ID)
		}
		return e.dbsOrdersRepo.MarkNotified(ctx, ids, now)
	}
	if len(products) == 0 {
		products, err = e.productsRepo.ListByShop(ctx, shopID, 0, 0)
		if err != nil {
			return err
		}
	}
	nmIDSet := make(map[int64]struct{}, len(orders))
	for _, o := range orders {
		nmIDSet[o.NmID] = struct{}{}
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
	text := notifications.FormatAssemblyOrdersMessageRU("DBS", dbsOrdersToAssemblyRows(orders), titleByNmID)
	now := time.Now()
	sent := false
	for _, ep := range endpoints {
		if ep.Channel != "telegram" || !ep.Enabled {
			continue
		}
		chatID, err := strconv.ParseInt(ep.Address, 10, 64)
		if err != nil {
			e.log.Warn("dbs orders: invalid telegram chat_id", "endpoint_id", ep.ID, "address", ep.Address, "err", err)
			continue
		}
		if err := e.tg.SendMessage(ctx, chatID, text); err != nil {
			e.log.Error("dbs orders: send failed", "endpoint_id", ep.ID, "err", err)
			continue
		}
		sent = true
	}
	if sent {
		ids := make([]int, 0, len(orders))
		for _, o := range orders {
			ids = append(ids, o.ID)
		}
		if err := e.dbsOrdersRepo.MarkNotified(ctx, ids, now); err != nil {
			return err
		}
	}
	return nil
}

func dbsOrdersToAssemblyRows(orders []*domain.DBSOrder) []notifications.AssemblyOrderRow {
	rows := make([]notifications.AssemblyOrderRow, 0, len(orders))
	for _, o := range orders {
		rows = append(rows, notifications.AssemblyOrderRow{
			NmID:      o.NmID,
			Article:   o.Article,
			WBOrderID: o.WBOrderID,
			CreatedAt: o.CreatedAt,
		})
	}
	return rows
}
