package dto

import "time"

// GetNewDBWOrdersResponse — ответ GET /api/v3/dbw/orders/new (сборочные задания DBW).
type GetNewDBWOrdersResponse struct {
	Orders []DBWOrderNew `json:"orders"`
}

// DBWOrderNew — одно новое сборочное задание DBW из WB API.
type DBWOrderNew struct {
	ID               int64          `json:"id"`
	OrderUID         string         `json:"orderUid"`
	GroupID          string         `json:"groupId"`
	NmID             int64          `json:"nmId"`
	ChrtID           int64          `json:"chrtId"`
	Article          string         `json:"article"`
	ColorCode        string         `json:"colorCode"`
	Rid              string         `json:"rid"`
	CreatedAt        time.Time      `json:"createdAt"`
	Skus             []string       `json:"skus"`
	WarehouseID      int64          `json:"warehouseId"`
	Price            int            `json:"price"`
	ConvertedPrice   int            `json:"convertedPrice"`
	CurrencyCode     int            `json:"currencyCode"`
	ConvertedCurrCode int           `json:"convertedCurrencyCode"`
	CargoType        int            `json:"cargoType"`
	IsZeroOrder      bool           `json:"isZeroOrder"`
	Comment          string         `json:"comment"`
	RequiredMeta     []string       `json:"requiredMeta"`
	Address          *DBWAddress    `json:"address,omitempty"`
	Options          *DBWOptions    `json:"options,omitempty"`
	SalePrice        int            `json:"salePrice,omitempty"`
}

// DBWAddress — адрес доставки DBW.
type DBWAddress struct {
	FullAddress string  `json:"fullAddress"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
}

// DBWOptions — опции заказа DBW.
type DBWOptions struct {
	IsB2B bool `json:"isB2b"`
}
