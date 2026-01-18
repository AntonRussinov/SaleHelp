package alert

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (e *AlertEngine) evalPriceDropEffective(
	ctx context.Context,
	products []*domain.Product,
	snapshotsByProduct map[int][]*domain.ProductSnapshot,
	global *domain.Alert,
	overrides map[int]*domain.Alert,
) error {

	var globalCfg priceDropCfg
	var hasGlobal bool

	if global != nil {
		cfg, err := parsePriceDropCfg(global.Payload)
		if err != nil {
			return fmt.Errorf("invalid global price_drop payload (alert_id=%d): %w", global.ID, err)
		}
		globalCfg = cfg
		hasGlobal = true
	}

	cfgCache := make(map[int]priceDropCfg) // alertID -> cfg

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

		snaps := snapshotsByProduct[p.ID]
		if len(snaps) < 2 {
			continue
		}

		latest := snaps[0]
		prev := snaps[1]

		latestPrice := derefFloat64(latest.FinalPrice)
		prevPrice := derefFloat64(prev.FinalPrice)
		// cfg
		var cfg priceDropCfg
		if useOverride {
			if cached, ok := cfgCache[a.ID]; ok {
				cfg = cached
			} else {
				parsed, err := parsePriceDropCfg(a.Payload)
				if err != nil {
					return fmt.Errorf("invalid price_drop payload (alert_id=%d): %w", a.ID, err)
				}
				cfgCache[a.ID] = parsed
				cfg = parsed
			}
		} else {
			if !hasGlobal {
				continue
			}
			cfg = globalCfg
		}

		ok, delta, percent := priceDropTriggered(prevPrice, latestPrice, cfg)
		if !ok {
			continue
		}

		msg := fmt.Sprintf(
			"Price dropped: %.2f → %.2f (-%.2f, %.2f%%)",
			prevPrice, latestPrice, delta, percent,
		)

		tr := domain.AlertTriggered{
			AlertID:    a.ID,
			ProductID:  p.ID,
			SnapshotID: latest.ID,
			ShopID:     p.ShopID,
			AlertType:  a.Type,
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

func priceDropTriggered(prevPrice, latestPrice float64, cfg priceDropCfg) (bool, float64, float64) {
	if prevPrice <= 0 || latestPrice <= 0 {
		return false, 0, 0
	}
	if latestPrice >= prevPrice {
		return false, 0, 0
	}

	delta := prevPrice - latestPrice
	percent := (float64(delta) / float64(prevPrice)) * 100.0

	okAmount := cfg.Amount > 0 && delta >= float64(cfg.Amount)
	okPercent := cfg.Percent > 0 && percent >= cfg.Percent

	var ok bool
	if cfg.Mode == "all" {
		// если задан только один порог — all эквивалентен ему
		switch {
		case cfg.Amount > 0 && cfg.Percent > 0:
			ok = okAmount && okPercent
		case cfg.Amount > 0:
			ok = okAmount
		default:
			ok = okPercent
		}
	} else { // any
		ok = okAmount || okPercent
	}

	return ok, delta, percent
}
