package alert

import (
	"context"
	"fmt"
	"strconv"

	"salehelp/domain"
)

func (e *AlertEngine) evalStockLowEffective(
	ctx context.Context,
	products []*domain.Product,
	sizesByProduct map[int][]*domain.ProductSize,
	snapshotsByProduct map[int][]*domain.ProductSnapshot,
	global *domain.Alert,
	overrides map[int]*domain.Alert,
	aProductIDs map[int]bool,
) error {

	// парсим конфиг глобального один раз
	var globalCfg stockLowCfg
	var hasGlobal bool

	if global != nil {
		cfg, err := parseStockLowCfg(global.Payload)
		if err != nil {
			return fmt.Errorf("invalid global stock_low payload (alert_id=%d): %w", global.ID, err)
		}
		globalCfg = cfg
		hasGlobal = true
	}

	// overrides парсим лениво (на товар) или можно заранее закэшировать
	overrideCfgCache := make(map[int]stockLowCfg)

	for _, p := range products {
		// выбираем effective alert
		a := global
		useOverride := false

		if ov, ok := overrides[p.ID]; ok && ov != nil {
			a = ov
			useOverride = true
		}

		// если нет ни global, ни override — нечего делать
		if a == nil {
			continue
		}

		hasSizes := len(sizesByProduct[p.ID]) > 0
		scope := chooseStockScope(a, hasSizes)
		if scope != "product" {
			continue
		}

		// снапшоты
		snaps := snapshotsByProduct[p.ID]
		if len(snaps) < 2 { // prev==nil -> не триггерим
			continue
		}

		latest := snaps[0]
		prev := snaps[1]

		latestStock := derefInt(latest.StockTotal)
		prevStock := derefInt(prev.StockTotal)

		if latestStock == 0 {
			continue
		}

		// cfg
		var cfg stockLowCfg
		if useOverride {
			if cached, ok := overrideCfgCache[a.ID]; ok {
				cfg = cached
			} else {
				parsed, err := parseStockLowCfg(a.Payload)
				if err != nil {
					return fmt.Errorf("invalid stock_low payload (alert_id=%d): %w", a.ID, err)
				}
				overrideCfgCache[a.ID] = parsed
				cfg = parsed
			}
		} else {
			// global
			if !hasGlobal {
				continue
			}
			cfg = globalCfg
		}

		effectiveThreshold := cfg.Threshold
		if aProductIDs[p.ID] {
			// A-товары: порог выше (раньше предупреждать)
			if effectiveThreshold < 10 {
				effectiveThreshold = 10
			}
			effectiveThreshold += 10
		}
		// edge: пересекли порог сверху вниз
		if !(prevStock > effectiveThreshold && latestStock <= effectiveThreshold) {
			continue
		}

		var msg string
		if aProductIDs[p.ID] {
			msg = "⚠️ A-товар заканчивается! " + p.Title + " — осталось " + strconv.Itoa(latestStock) + " шт"
		} else {
			msg = fmt.Sprintf("Stock is low: %d (threshold %d)", latestStock, effectiveThreshold)
		}

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

func (e *AlertEngine) evalStockLowSizesEffective(
	ctx context.Context,
	products []*domain.Product,
	sizesByProduct map[int][]*domain.ProductSize,
	sizeSnapsBySize map[int][]*domain.ProductSizeSnapshot,
	global *domain.Alert,
	overrides map[int]*domain.Alert,
	aProductIDs map[int]bool,
) error {

	var globalCfg stockLowCfg
	var hasGlobal bool
	if global != nil {
		cfg, err := parseStockLowCfg(global.Payload)
		if err != nil {
			return fmt.Errorf("invalid global stock_low payload (alert_id=%d): %w", global.ID, err)
		}
		globalCfg = cfg
		hasGlobal = true
	}

	overrideCfgCache := make(map[int]stockLowCfg)

	for _, p := range products {
		// effective alert
		a := global
		useOverride := false
		if ov, ok := overrides[p.ID]; ok && ov != nil {
			a = ov
			useOverride = true
		}
		if a == nil {
			continue
		}

		hasSizes := len(sizesByProduct[p.ID]) > 0
		scope := chooseStockScope(a, hasSizes)
		if scope != "size" {
			continue
		}

		// sizes for product
		sizes := sizesByProduct[p.ID]
		if len(sizes) == 0 {
			continue
		}

		// cfg
		var cfg stockLowCfg
		if useOverride {
			if cached, ok := overrideCfgCache[a.ID]; ok {
				cfg = cached
			} else {
				parsed, err := parseStockLowCfg(a.Payload)
				if err != nil {
					return fmt.Errorf("invalid stock_low payload (alert_id=%d): %w", a.ID, err)
				}
				overrideCfgCache[a.ID] = parsed
				cfg = parsed
			}
		} else {
			if !hasGlobal {
				continue
			}
			cfg = globalCfg
		}

		// optional: если задан alert.ProductSizeID — проверяем только этот размер
		for _, sz := range sizes {
			if a.ProductSizeID != nil && *a.ProductSizeID != sz.ID {
				continue
			}

			snaps := sizeSnapsBySize[sz.ID]
			if len(snaps) < 2 {
				continue // prev==nil -> не триггерим edge
			}

			latest := snaps[0]
			prev := snaps[1]

			latestStock := derefInt(latest.Stock)
			prevStock := derefInt(prev.Stock)

			// priority: if stock is zero -> only stock_zero should fire
			if latestStock == 0 {
				continue
			}

			effectiveThreshold := cfg.Threshold
			if aProductIDs[p.ID] {
				if effectiveThreshold < 10 {
					effectiveThreshold = 10
				}
				effectiveThreshold += 10
			}
			if !(prevStock > effectiveThreshold && latestStock <= effectiveThreshold) {
				continue
			}

			var msg string
			if aProductIDs[p.ID] {
				msg = "⚠️ A-товар заканчивается! " + p.Title + " — осталось " + strconv.Itoa(latestStock) + " шт"
			} else {
				msg = fmt.Sprintf("Низкий остаток: %d шт (порог %d)", latestStock, effectiveThreshold)
			}

			tr := domain.AlertTriggered{
				AlertID:       a.ID,
				ShopID:        p.ShopID,
				ProductID:     p.ID,
				ProductSizeID: &sz.ID,                   // <-- важно
				SnapshotID:    latest.ProductSnapshotID, // <-- смотри ниже примечание
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

			// fanout/deliveries создавай по tr.ID (как у тебя уже сделано)

			// ВОТ ТУТ fanout
			if err := e.deliveriesFanoutRepo.CreateDeliveriesForTriggered(ctx, ID, tr.ShopID); err != nil {
				return err
			}
		}
	}

	return nil
}
