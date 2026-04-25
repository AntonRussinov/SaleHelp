package alert

import (
	"context"
)

// buildStockWarningByNmID возвращает множество nmID, у которых количество активных FBS-заказов
// (supplier_status='new') >= остатку из последнего снапшота. Вызывается при формировании
// уведомления о заказе, чтобы сразу пометить «последний товар» без ожидания следующего синка.
func (e *AlertEngine) buildStockWarningByNmID(
	ctx context.Context,
	shopID int,
	nmIDSet map[int64]struct{},
	productIDByNmID map[int64]int,
) map[int64]bool {
	result := make(map[int64]bool)
	if len(nmIDSet) == 0 || len(productIDByNmID) == 0 {
		return result
	}

	// Pending FBS-заказы по nmID (supplier_status='new').
	pendingByNmID, err := e.fbsOrdersRepo.CountNewByShopGroupedByNmID(ctx, shopID)
	if err != nil {
		e.log.Warn("stock_warning: count new fbs orders failed", "shop_id", shopID, "err", err)
		return result
	}
	if len(pendingByNmID) == 0 {
		return result
	}

	// Последние снапшоты только по продуктам из текущих заказов.
	productIDs := make([]int, 0, len(productIDByNmID))
	for _, pid := range productIDByNmID {
		productIDs = append(productIDs, pid)
	}

	snapshotsByProduct, err := e.snapshotsRepo.GetLatestByProductIDsLimit(ctx, productIDs, 1)
	if err != nil {
		e.log.Warn("stock_warning: get snapshots failed", "shop_id", shopID, "err", err)
		return result
	}

	for nmID := range nmIDSet {
		pending, ok := pendingByNmID[nmID]
		if !ok || pending == 0 {
			continue
		}
		productID, ok := productIDByNmID[nmID]
		if !ok {
			continue
		}
		snaps := snapshotsByProduct[productID]
		if len(snaps) == 0 || snaps[0].StockTotal == nil {
			continue
		}
		stock := *snaps[0].StockTotal
		if stock > 0 && pending >= stock {
			result[nmID] = true
		}
	}

	return result
}

