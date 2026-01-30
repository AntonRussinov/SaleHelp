package alert

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/services/subscription"
)

func (e *AlertEngine) RunShop(ctx context.Context, shopID int) error {

	lvl, _, err := e.subSvc.GetShopAccessLevel(ctx, shopID)
	if err != nil {
		return err
	}

	if lvl != subscription.AccessActive {
		return nil // алерты отключены (grace/blocked/none)
	}

	alerts, err := e.alertsRepo.ListActiveByShop(ctx, shopID)
	if err != nil {
		return err
	}

	var products []*domain.Product
	if len(alerts) > 0 {
		products, err = e.productsRepo.ListByShop(ctx, shopID, 0, 0)
		if err != nil || len(products) == 0 {
			return err
		}
	}

	// batch: 1 запрос на latest+prev (только если есть продукты для алертов)
	if len(products) > 0 {
		productIDs := make([]int, 0, len(products))
		for _, p := range products {
			productIDs = append(productIDs, p.ID)
		}

		sizes, err := e.productSizesRepo.ListByProductIDs(ctx, productIDs)
		if err != nil {
			return err
		}

		sizesByProduct := make(map[int][]*domain.ProductSize, len(products))
		for _, s := range sizes {
			sizesByProduct[s.ProductID] = append(sizesByProduct[s.ProductID], s)
		}

		sizeIDs := make([]int, 0, len(sizes))
		for _, s := range sizes {
			sizeIDs = append(sizeIDs, s.ID)
		}

		sizeSnapsBySize := map[int][]*domain.ProductSizeSnapshot{}
		if len(sizeIDs) > 0 {
			sizeSnapsBySize, err = e.productSizeSnapshotsRepo.GetLatestBySizeIDsLimitMap(ctx, sizeIDs, 2)
			if err != nil {
				return err
			}
		}

		snapshotsByProduct, err := e.snapshotsRepo.GetLatestByProductIDsLimit(ctx, productIDs, 2)
		if err != nil {
			return err
		}

		// --- stock_low: global + overrides ---
		globalLow, lowOverrides := buildEffectiveAlertsByProduct(alerts, domain.AlertTypeStockLow)
		if globalLow != nil || len(lowOverrides) > 0 {
			if err := e.evalStockLowEffective(ctx, products, sizesByProduct, snapshotsByProduct, globalLow, lowOverrides); err != nil {
				return err
			}

			if err := e.evalStockLowSizesEffective(ctx, products, sizesByProduct, sizeSnapsBySize, globalLow, lowOverrides); err != nil {
				return err
			}
		}

		// --- stock_zero: global + overrides ---
		globalZero, zeroOverrides := buildEffectiveAlertsByProduct(alerts, domain.AlertTypeStockZero)
		if globalZero != nil || len(zeroOverrides) > 0 {
			if err := e.evalStockZeroEffective(ctx, products, sizesByProduct, snapshotsByProduct, globalZero, zeroOverrides); err != nil {
				return err
			}

			if err := e.evalStockZeroSizesEffective(ctx, products, sizesByProduct, sizeSnapsBySize, globalZero, zeroOverrides); err != nil {
				return err
			}
		}

		globalDrop, dropOverrides := buildEffectiveAlertsByProduct(alerts, domain.AlertTypePriceDrop)
		if globalDrop != nil || len(dropOverrides) > 0 {
			if err := e.evalPriceDropEffective(ctx, products, snapshotsByProduct, globalDrop, dropOverrides); err != nil {
				return err
			}
		}
	}

	// новые FBS-заказы: отправить в Telegram и пометить как уведомлённые (products — для подстановки названий; при nil подгрузим внутри)
	if err := e.runShopFBSOrders(ctx, shopID, products); err != nil {
		return err
	}

	return nil
}
