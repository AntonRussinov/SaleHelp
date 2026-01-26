package alert

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (e *AlertEngine) evalStockZeroEffective(
	ctx context.Context,
	products []*domain.Product,
	sizesByProduct map[int][]*domain.ProductSize,
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

		hasSizes := len(sizesByProduct[p.ID]) > 0
		scope := chooseStockScope(a, hasSizes)
		if scope != "product" {
			continue
		}

		snaps := snapshotsByProduct[p.ID]
		if len(snaps) < 2 { // prev==nil -> не триггерим
			continue
		}

		latest := snaps[0]
		prev := snaps[1]

		// Если WB stocks не был получен, StockTotal может быть nil (unknown).
		// В этом случае нельзя трактовать nil как 0 и триггерить "всё закончилось".
		if latest.StockTotal == nil || prev.StockTotal == nil {
			continue
		}

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
			AlertType:  a.Type,
			Message:    msg,
		}

		ID, inserted, err := e.alertsTriggeredRepo.CreateIfNotExists(ctx, &tr)
		if err != nil {
			return err
		}

		if !inserted {
			continue
		}

		// ВОТ ТУТ fanout
		if err := e.deliveriesFanoutRepo.CreateDeliveriesForTriggered(ctx, ID, tr.ShopID); err != nil {
			return err
		}

	}

	return nil
}

func (e *AlertEngine) evalStockZeroSizesEffective(
	ctx context.Context,
	products []*domain.Product,
	sizesByProduct map[int][]*domain.ProductSize,
	sizeSnapsBySize map[int][]*domain.ProductSizeSnapshot,
	global *domain.Alert,
	overrides map[int]*domain.Alert,
) error {
	for _, p := range products {
		a := global
		if ov, ok := overrides[p.ID]; ok && ov != nil {
			a = ov
		}
		if a == nil {
			continue
		}

		hasSizes := len(sizesByProduct[p.ID]) > 0
		scope := chooseStockScope(a, hasSizes)
		if scope != "size" {
			continue
		}

		sizes := sizesByProduct[p.ID]
		if len(sizes) == 0 {
			continue
		}

		for _, sz := range sizes {
			if a.ProductSizeID != nil && *a.ProductSizeID != sz.ID {
				continue
			}

			snaps := sizeSnapsBySize[sz.ID]
			if len(snaps) < 2 {
				continue
			}

			latest := snaps[0]
			prev := snaps[1]

			// Аналогично по размеру: nil = нет данных, не 0.
			if latest.Stock == nil || prev.Stock == nil {
				continue
			}

			latestStock := derefInt(latest.Stock)
			prevStock := derefInt(prev.Stock)

			if !(prevStock > 0 && latestStock == 0) {
				continue
			}

			msg := fmt.Sprintf("Размер закончился: %d → %d", prevStock, latestStock)

			tr := domain.AlertTriggered{
				AlertID:       a.ID,
				ShopID:        p.ShopID,
				ProductID:     p.ID,
				ProductSizeID: &sz.ID,
				SnapshotID:    latest.ProductSnapshotID,
				AlertType:     a.Type,
				Message:       msg,
			}

			ID, inserted, err := e.alertsTriggeredRepo.CreateIfNotExists(ctx, &tr)
			if err != nil {
				return err
			}
			if !inserted {
				continue
			}

			// ВОТ ТУТ fanout
			if err := e.deliveriesFanoutRepo.CreateDeliveriesForTriggered(ctx, ID, tr.ShopID); err != nil {
				return err
			}
		}
	}

	return nil
}
