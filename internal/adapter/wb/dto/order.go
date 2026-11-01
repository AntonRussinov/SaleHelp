package dto

import "time"

// OrderItem — одна запись из отчёта заказов WB (GET /api/v1/supplier/orders).
// Документация: https://dev.wildberries.ru/docs/openapi/reports#tag/Osnovnye-otchyoty/paths/~1api~1v1~1supplier~1orders/get
// Связь с Sales по полю srid. isCancel = true — заказ не забран (отказ).
type OrderItem struct {
	Date            string  `json:"date"`
	LastChangeDate  string  `json:"lastChangeDate"`
	WarehouseName   string  `json:"warehouseName"`
	WarehouseType   string  `json:"warehouseType"`
	CountryName     string  `json:"countryName"`
	OblastOkrugName string  `json:"oblastOkrugName"`
	RegionName      string  `json:"regionName"`
	SupplierArticle string  `json:"supplierArticle"`
	NmID            int64   `json:"nmId"`
	Barcode         string  `json:"barcode"`
	Category        string  `json:"category"`
	Subject         string  `json:"subject"`
	Brand           string  `json:"brand"`
	TechSize        string  `json:"techSize"`
	IncomeID        int64   `json:"incomeID"`
	IsSupply        bool    `json:"isSupply"`
	IsRealization   bool    `json:"isRealization"`
	TotalPrice      float64 `json:"totalPrice"`
	DiscountPercent float64 `json:"discountPercent"`
	Spp             float64 `json:"spp"`
	FinishedPrice   float64 `json:"finishedPrice"`
	PriceWithDisc   float64 `json:"priceWithDisc"`
	IsCancel        bool    `json:"isCancel"`
	CancelDate      string  `json:"cancelDate"`
	Sticker         string  `json:"sticker"`
	GNumber         string  `json:"gNumber"`
	Srid            string  `json:"srid"`
}

// GetOrdersResponse — ответ GET /api/v1/supplier/orders (массив записей).
type GetOrdersResponse []OrderItem

// ParseDate парсит дату из WB. loc — пояс магазина; nil = глобальный из конфига.
func (o *OrderItem) ParseDate(loc *time.Location) (time.Time, error) {
	return parseWBTime(o.Date, loc)
}

// ParseLastChangeDate парсит lastChangeDate из WB. loc — пояс магазина; nil = глобальный.
func (o *OrderItem) ParseLastChangeDate(loc *time.Location) (time.Time, error) {
	return parseWBTime(o.LastChangeDate, loc)
}

// ParseCancelDate парсит cancelDate из WB. loc — пояс магазина; nil = глобальный.
func (o *OrderItem) ParseCancelDate(loc *time.Location) (time.Time, error) {
	return parseWBTime(o.CancelDate, loc)
}
