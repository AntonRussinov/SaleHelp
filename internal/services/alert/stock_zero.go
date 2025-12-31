package alert

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (e *AlertEngine) evalStockZeroEffective(
	ctx context.Context,
	products []*domain.Product,
	snapshotsByProduct map[int][]*domain.ProductSnapshot,
	global *domain.Alert,
	overrides map[int]*domain.Alert,
) error {

	for _, p := range products {
		// effective alert
		a := global
		if ov, ok := overrides[p.ID]; ok && ov != nil {
			a = ov
		}
		if a == nil {
			continue
		}

		snaps := snapshotsByProduct[p.ID]
		if len(snaps) < 2 { // prev==nil -> не триггерим
			continue
		}

		latest := snaps[0]
		prev := snaps[1]

		latestStock := derefInt(latest.StockTotal)
		prevStock := derefInt(prev.StockTotal)

		// edge: >0 -> 0
		if !(prevStock > 0 && latestStock == 0) {
			continue
		}

		msg := fmt.Sprintf("Stock is ZERO: %d → %d", prevStock, latestStock)

		tr := domain.AlertTriggered{
			AlertID:    a.ID,
			ProductID:  p.ID,
			SnapshotID: latest.ID,
			ShopID:     p.ShopID,
			Message:    msg,
		}

		id, inserted, err := e.alertsTriggeredRepo.CreateIfNotExists(ctx, &tr)
		if err != nil {
			return err
		}

		if !inserted {
			continue
		}

		// ВОТ ТУТ fanout
		if err := e.deliveriesFanoutRepo.CreateDeliveriesForTriggered(ctx, id, tr.ShopID); err != nil {
			return err
		}

	}

	return nil
}
