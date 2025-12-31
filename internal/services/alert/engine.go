package alert

import (
	"context"
)

func (e *AlertEngine) RunShop(ctx context.Context, shopID int) error {
	alerts, err := e.alertsRepo.ListActiveByShop(ctx, shopID)
	if err != nil || len(alerts) == 0 {
		return err
	}

	products, err := e.productsRepo.ListByShop(ctx, shopID, 0, 0)
	if err != nil || len(products) == 0 {
		return err
	}

	// batch: 1 запрос на latest+prev
	productIDs := make([]int, 0, len(products))
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}

	snapshotsByProduct, err := e.snapshotsRepo.GetLatestByProductIDsLimit(ctx, productIDs, 2)
	if err != nil {
		return err
	}

	// --- stock_low: global + overrides ---
	globalLow, lowOverrides := buildEffectiveAlertsByProduct(alerts, AlertTypeStockLow)
	if globalLow != nil || len(lowOverrides) > 0 {
		if err := e.evalStockLowEffective(ctx, products, snapshotsByProduct, globalLow, lowOverrides); err != nil {
			return err
		}
	}

	// --- stock_zero: global + overrides ---
	globalZero, zeroOverrides := buildEffectiveAlertsByProduct(alerts, AlertTypeStockZero)
	if globalZero != nil || len(zeroOverrides) > 0 {
		if err := e.evalStockZeroEffective(ctx, products, snapshotsByProduct, globalZero, zeroOverrides); err != nil {
			return err
		}
	}

	globalDrop, dropOverrides := buildEffectiveAlertsByProduct(alerts, AlertTypePriceDrop)
	if globalDrop != nil || len(dropOverrides) > 0 {
		if err := e.evalPriceDropEffective(ctx, products, snapshotsByProduct, globalDrop, dropOverrides); err != nil {
			return err
		}
	}

	return nil
}
