package dto

import (
	"time"
)

// Часовой пояс WB для отчётов. Задаётся при старте приложения из конфига (WBReportTimezone).
// Если не задан — Europe/Moscow.
var wbReportLoc *time.Location

func init() {
	var err error
	wbReportLoc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		wbReportLoc = time.FixedZone("MSK", 3*3600)
	}
}

// SetWBReportLocation задаёт часовой пояс для парсинга дат из отчётов WB (sales, orders).
// Вызывается при старте из конфига. Например Asia/Almaty, если WB отдаёт время по времени продавца.
func SetWBReportLocation(loc *time.Location) {
	if loc != nil {
		wbReportLoc = loc
	}
}

// SaleItem — одна запись из отчёта продаж/возвратов WB (GET /api/v1/supplier/sales).
// Документация: https://dev.wildberries.ru/openapi/reports#tag/Osnovnye-otchyoty/paths/~1api~1v1~1supplier~1sales/get
type SaleItem struct {
	Date              string  `json:"date"`
	LastChangeDate    string  `json:"lastChangeDate"`
	WarehouseName     string  `json:"warehouseName"`
	WarehouseType     string  `json:"warehouseType"`
	CountryName       string  `json:"countryName"`
	OblastOkrugName   string  `json:"oblastOkrugName"`
	RegionName        string  `json:"regionName"`
	SupplierArticle   string  `json:"supplierArticle"`
	NmID              int64   `json:"nmId"`
	Barcode           string  `json:"barcode"`
	Category          string  `json:"category"`
	Subject           string  `json:"subject"`
	Brand             string  `json:"brand"`
	TechSize          string  `json:"techSize"`
	IncomeID          int64   `json:"incomeID"`
	IsSupply          bool    `json:"isSupply"`
	IsRealization     bool    `json:"isRealization"`
	TotalPrice        float64 `json:"totalPrice"`
	DiscountPercent   float64 `json:"discountPercent"`
	Spp               float64 `json:"spp"`
	PaymentSaleAmount float64 `json:"paymentSaleAmount"`
	ForPay            float64 `json:"forPay"`
	FinishedPrice     float64 `json:"finishedPrice"`
	PriceWithDisc     float64 `json:"priceWithDisc"`
	SaleID            string  `json:"saleID"`
	Sticker           string  `json:"sticker"`
	GNumber           string  `json:"gNumber"`
	Srid              string  `json:"srid"`
}

// GetSalesResponse — ответ GET /api/v1/supplier/sales (массив записей).
type GetSalesResponse []SaleItem

// parseWBTime парсит дату из WB. Если в строке нет зоны (Z, +03:00) — используется loc (или wbReportLoc при loc == nil).
func parseWBTime(s string, loc *time.Location) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	use := wbReportLoc
	if loc != nil {
		use = loc
	}
	t, err := time.ParseInLocation("2006-01-02T15:04:05", s, use)
	if err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02T15:04:05.9", s, use)
}

// ParseDate парсит дату из WB. loc — пояс магазина для отчётов; nil = глобальный из конфига.
func (s *SaleItem) ParseDate(loc *time.Location) (time.Time, error) {
	return parseWBTime(s.Date, loc)
}

// ParseLastChangeDate парсит lastChangeDate из WB. loc — пояс магазина; nil = глобальный.
func (s *SaleItem) ParseLastChangeDate(loc *time.Location) (time.Time, error) {
	return parseWBTime(s.LastChangeDate, loc)
}
