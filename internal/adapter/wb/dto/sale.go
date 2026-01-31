package dto

import "time"

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

// parseWBTime парсит дату из WB (RFC3339 или "2006-01-02T15:04:05").
func parseWBTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02T15:04:05", s)
}

// ParseDate парсит дату из WB.
func (s *SaleItem) ParseDate() (time.Time, error) {
	return parseWBTime(s.Date)
}

// ParseLastChangeDate парсит lastChangeDate из WB.
func (s *SaleItem) ParseLastChangeDate() (time.Time, error) {
	return parseWBTime(s.LastChangeDate)
}
