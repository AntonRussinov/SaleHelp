package monitoring

import "salehelp/domain"

// LatestProductState — “текущая картинка” товара на основе последнего snapshot.
type LatestProductState struct {
	ProductID int
	Snapshot  *domain.ProductSnapshot

	// Stock: (product_size_id, sku) -> qty
	StockBySizeSKU map[SizeSKUKey]int

	// Stock: product_size_id -> sum(qty by sku)
	StockBySizeID map[int]int

	// Prices: product_size_id -> price snapshot (final/base/etc)
	PriceBySizeID map[int]*domain.ProductSizePriceSnapshot
}

type SizeSKUKey struct {
	ProductSizeID int
	SKU           string
}
