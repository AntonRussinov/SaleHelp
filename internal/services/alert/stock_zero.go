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
		// Нужно минимум 3 снэпшота: текущий + предыдущий + позапрошлый.
		// Два consecutive нуля подряд после ненулевого — подтверждение реального обнуления.
		// Это защищает от ложных алертов при временном сбое WB API (isolated zero).
		if len(snaps) < 3 {
			continue
		}

		latest := snaps[0]
		prev := snaps[1]
		prevPrev := snaps[2]

		// nil = данные не были получены от WB, не трактуем как 0.
		if latest.StockTotal == nil || prev.StockTotal == nil || prevPrev.StockTotal == nil {
			continue
		}

		latestStock := derefInt(latest.StockTotal)
		prevStock := derefInt(prev.StockTotal)
		prevPrevStock := derefInt(prevPrev.StockTotal)

		// Требуем: два нуля подряд после ненулевого значения.
		// latest=0, prev=0, prevPrev>0 → первый "подтверждённый" тик обнуления.
		if !(latestStock == 0 && prevStock == 0 && prevPrevStock > 0) {
			continue
		}

		msg := fmt.Sprintf("Stock is ZERO: %d → %d → %d", prevPrevStock, prevStock, latestStock)

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
			if len(snaps) < 3 {
				continue
			}

			latest := snaps[0]
			prev := snaps[1]
			prevPrev := snaps[2]

			if latest.Stock == nil || prev.Stock == nil || prevPrev.Stock == nil {
				continue
			}

			latestStock := derefInt(latest.Stock)
			prevStock := derefInt(prev.Stock)
			prevPrevStock := derefInt(prevPrev.Stock)

			if !(latestStock == 0 && prevStock == 0 && prevPrevStock > 0) {
				continue
			}

			msg := fmt.Sprintf("Размер закончился: %d → %d → %d", prevPrevStock, prevStock, latestStock)

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

			if err := e.deliveriesFanoutRepo.CreateDeliveriesForTriggered(ctx, ID, tr.ShopID); err != nil {
				return err
			}
		}
	}

	return nil
}
