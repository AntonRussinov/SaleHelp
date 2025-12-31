package domain

type ProductSizeSnapshot struct {
	ID                int
	ProductSnapshotID int
	ProductSizeID     int
	SKU               string

	Stock            *int
	StockByWarehouse []byte
	Extra            []byte
}
