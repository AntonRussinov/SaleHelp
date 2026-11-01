package dto

import "time"

// GetNewDBSOrdersResponse — ответ GET /api/v3/dbs/orders/new (сборочные задания DBS).
type GetNewDBSOrdersResponse struct {
	Orders []DBSOrderNew `json:"orders"`
}

// DBSOrderNew — одно новое сборочное задание DBS из WB API.
type DBSOrderNew struct {
	ID                 int64          `json:"id"`
	OrderUID           string         `json:"orderUid"`
	GroupID            string         `json:"groupId"`
	NmID               int64          `json:"nmId"`
	ChrtID             int64          `json:"chrtId"`
	Article            string         `json:"article"`
	ColorCode          string         `json:"colorCode"`
	Rid                string         `json:"rid"`
	CreatedAt          time.Time      `json:"createdAt"`
	Skus               []string       `json:"skus"`
	WarehouseID        int64          `json:"warehouseId"`
	Price              int            `json:"price"`
	FinalPrice         int            `json:"finalPrice"`
	ConvertedPrice     int            `json:"convertedPrice"`
	ConvertedFinalPrice int           `json:"convertedFinalPrice"`
	CurrencyCode       int            `json:"currencyCode"`
	ConvertedCurrCode  int            `json:"convertedCurrencyCode"`
	CargoType          int            `json:"cargoType"`
	IsZeroOrder        bool           `json:"isZeroOrder"`
	Comment            string         `json:"comment"`
	RequiredMeta       []string       `json:"requiredMeta"`
	DeliveryType       string         `json:"deliveryType"`
	Address            *DBSAddress    `json:"address,omitempty"`
	Options            *DBSOptions    `json:"options,omitempty"`
	SalePrice          int            `json:"salePrice,omitempty"`
}

// DBSAddress — адрес доставки DBS.
type DBSAddress struct {
	FullAddress string  `json:"fullAddress"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
}

// DBSOptions — опции заказа DBS.
type DBSOptions struct {
	IsB2B bool `json:"isB2b"`
}
