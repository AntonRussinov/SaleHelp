package alert

import (
	"context"

	"salehelp/domain"
)

// buildStockWarningByNmID возвращает множество nmID, у которых количество заказов в текущем
// уведомительном батче >= остатку из последнего снапшота. Вызывается при формировании
// уведомления о заказе, чтобы сразу пометить «последний товар» без ожидания следующего синка.
func (e *AlertEngine) buildStockWarningByNmID(
	ctx context.Context,
	shopID int,
	orders []*domain.FBSOrder,
	productIDByNmID map[int64]int,
) map[int64]bool {
	result := make(map[int64]bool)
	if len(orders) == 0 || len(productIDByNmID) == 0 {
		return result
	}

	// Считаем заказы из текущего батча (unnotified) — это реальный pending, а не все исторические.
	pendingByNmID := make(map[int64]int, len(orders))
	for _, o := range orders {
		pendingByNmID[o.NmID]++
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

	for nmID, pending := range pendingByNmID {
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

