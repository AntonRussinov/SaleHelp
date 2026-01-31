package domain

import "time"

// SaleEventType — тип события из отчёта продаж WB (выкуп, возврат, частичный возврат).
type SaleEventType string

const (
	SaleEventTypeBuyout         SaleEventType = "buyout"
	SaleEventTypeReturn         SaleEventType = "return"
	SaleEventTypePartialReturn SaleEventType = "partial_return"
)

// SaleEvent — запись из отчёта продаж/возвратов WB (GET /api/v1/supplier/sales).
// Используется для уведомлений о выкупах и возвратах в Telegram.
type SaleEvent struct {
	ID             int
	ShopID         int
	EventType      SaleEventType // buyout | return | partial_return
	WBSaleID       string  // saleID из WB (уникален в рамках отчёта)
	Srid           string  // srid — идентификатор заказа в WB
	Date           time.Time
	LastChangeDate time.Time
	NmID           int64
	Barcode        string
	SupplierArticle string
	WarehouseName  string
	Category       string
	Subject        string
	Brand          string
	TechSize       string
	IncomeID       int64
	TotalPrice     float64
	DiscountPercent float64
	FinishedPrice  float64
	PriceWithDisc  float64
	ForPay         float64
	CreatedAt      time.Time
	NotifiedAt     *time.Time
	Extra          []byte // JSON: regionName, countryName, spp, paymentSaleAmount и т.д.
}
