package dto

// SupplierStockItem — одна запись остатков со складов WB (GET /api/v1/supplier/stocks).
// Содержит текущий остаток товара на складах WB (включая возвраты, которые WB продаёт сам).
type SupplierStockItem struct {
	LastChangeDate  string `json:"lastChangeDate"`
	SupplierArticle string `json:"supplierArticle"`
	TechSize        string `json:"techSize"`
	Barcode         string `json:"barcode"`
	Quantity        int    `json:"quantity"`
	QuantityFull    int    `json:"quantityFull"`
	InWayToClient   int    `json:"inWayToClient"`
	InWayFromClient int    `json:"inWayFromClient"`
	WarehouseName   string `json:"warehouseName"`
	NmID            int64  `json:"nmId"`
	Subject         string `json:"subject"`
	Category        string `json:"category"`
	Brand           string `json:"brand"`
}

// GetSupplierStocksResponse — ответ GET /api/v1/supplier/stocks (массив записей).
type GetSupplierStocksResponse []SupplierStockItem
