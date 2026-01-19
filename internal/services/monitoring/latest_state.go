package monitoring

import (
	"context"
	"fmt"

	"salehelp/domain"
)

func (s *Service) GetLatestStateByProductID(ctx context.Context, productID int) (*LatestProductState, error) {
	// 1) последний snapshot
	snap, err := s.snapshotsRepo.GetLatestByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("get latest state: get latest snapshot: %w", err)
	}
	if snap == nil {
		// товар есть, но снапшота ещё не было
		return &LatestProductState{
			ProductID:      productID,
			Snapshot:       nil,
			StockBySizeSKU: map[SizeSKUKey]int{},
			StockBySizeID:  map[int]int{},
			PriceBySizeID:  map[int]*domain.ProductSizePriceSnapshot{},
		}, nil
	}

	// 2) stocks (product_size_snapshots) по snapshot_id
	stockRows, err := s.sizeSnapshotsRepo.ListBySnapshotID(ctx, snap.ID)
	if err != nil {
		return nil, fmt.Errorf("get latest state: list size stock snapshots: %w", err)
	}

	stockBySizeSKU := make(map[SizeSKUKey]int, len(stockRows))
	stockBySizeID := make(map[int]int)

	for _, r := range stockRows {
		if r == nil {
			continue
		}

		// Подстрой под твой domain.ProductSizeSnapshot:
		// ожидаю: ProductSizeID int, SKU string, Stock *int
		qty := 0
		if r.Stock != nil {
			qty = *r.Stock
		}

		key := SizeSKUKey{
			ProductSizeID: r.ProductSizeID,
			SKU:           r.SKU,
		}
		// если вдруг дубликаты — суммируем
		stockBySizeSKU[key] += qty
		stockBySizeID[r.ProductSizeID] += qty
	}

	// 3) prices (product_size_price_snapshots) по snapshot_id
	priceRows, err := s.sizePriceSnapshotsRepo.ListBySnapshotID(ctx, snap.ID)
	if err != nil {
		return nil, fmt.Errorf("get latest state: list size price snapshots: %w", err)
	}

	priceBySizeID := make(map[int]*domain.ProductSizePriceSnapshot, len(priceRows))
	for _, r := range priceRows {
		if r == nil {
			continue
		}
		// если вдруг несколько — оставляем “последний из выборки” (но по UNIQUE их быть не должно)
		priceBySizeID[r.ProductSizeID] = r
	}

	return &LatestProductState{
		ProductID:      productID,
		Snapshot:       snap,
		StockBySizeSKU: stockBySizeSKU,
		StockBySizeID:  stockBySizeID,
		PriceBySizeID:  priceBySizeID,
	}, nil
}

func (s *Service) GetLatestStateByProductIDs(ctx context.Context, productIDs []int) (map[int]*LatestProductState, error) {
	out := make(map[int]*LatestProductState, len(productIDs))
	if len(productIDs) == 0 {
		return out, nil
	}

	// 0) заранее создаём “пустые” стейты, чтобы продукт всегда был в ответе
	for _, pid := range productIDs {
		out[pid] = &LatestProductState{
			ProductID:      pid,
			Snapshot:       nil,
			StockBySizeSKU: map[SizeSKUKey]int{},
			StockBySizeID:  map[int]int{},
			PriceBySizeID:  map[int]*domain.ProductSizePriceSnapshot{},
		}
	}

	// 1) последние snapshots по продуктам (одним запросом)
	latestByProductID, err := s.snapshotsRepo.GetLatestByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("get latest state batch: get latest snapshots: %w", err)
	}

	// 2) соберём snapshotIDs и мапу snapshot_id -> product_id
	snapshotIDs := make([]int, 0, len(latestByProductID))
	snapIDToProductID := make(map[int]int, len(latestByProductID))

	for pid, snap := range latestByProductID {
		if snap == nil {
			continue
		}
		out[pid].Snapshot = snap
		snapshotIDs = append(snapshotIDs, snap.ID)
		snapIDToProductID[snap.ID] = pid
	}

	// если ни у одного продукта нет снапшота — возвращаем “пустые”
	if len(snapshotIDs) == 0 {
		return out, nil
	}

	// 3) stocks: все size stock snapshots по snapshotIDs (одним запросом)
	stockRows, err := s.sizeSnapshotsRepo.ListBySnapshotIDs(ctx, snapshotIDs)
	if err != nil {
		return nil, fmt.Errorf("get latest state batch: list size stock snapshots: %w", err)
	}

	for _, r := range stockRows {
		if r == nil {
			continue
		}

		pid, ok := snapIDToProductID[r.ProductSnapshotID]
		if !ok {
			continue
		}

		st := out[pid]
		if st == nil {
			continue
		}

		qty := 0
		if r.Stock != nil {
			qty = *r.Stock
		}

		key := SizeSKUKey{
			ProductSizeID: r.ProductSizeID,
			SKU:           r.SKU,
		}

		st.StockBySizeSKU[key] += qty
		st.StockBySizeID[r.ProductSizeID] += qty
	}

	// 4) prices: все size price snapshots по snapshotIDs (одним запросом)
	priceRows, err := s.sizePriceSnapshotsRepo.ListBySnapshotIDs(ctx, snapshotIDs)
	if err != nil {
		return nil, fmt.Errorf("get latest state batch: list size price snapshots: %w", err)
	}

	for _, r := range priceRows {
		if r == nil {
			continue
		}

		pid, ok := snapIDToProductID[r.ProductSnapshotID]
		if !ok {
			continue
		}

		st := out[pid]
		if st == nil {
			continue
		}

		// UNIQUE(product_snapshot_id, product_size_id) => дублей быть не должно
		st.PriceBySizeID[r.ProductSizeID] = r
	}

	return out, nil
}
