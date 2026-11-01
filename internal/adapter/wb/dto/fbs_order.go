package dto

import "time"

// GetNewOrdersResponse — ответ GET /api/v3/orders/new (сборочные задания FBS).
type GetNewOrdersResponse struct {
	Orders []FBSOrderNew `json:"orders"`
}

// FBSOrderNew — одно новое сборочное задание из WB API.
// Статусы supplierStatus/wbStatus приходят из POST /api/v3/orders/status.
type FBSOrderNew struct {
	ID             int64     `json:"id"`
	OrderUID       string    `json:"orderUid"`
	NmID           int64     `json:"nmId"`
	ChrtID         int64     `json:"chrtId"`
	Article        string    `json:"article"`
	CreatedAt      time.Time `json:"createdAt"`
	DeliveryType   string    `json:"deliveryType"`
	Price          int       `json:"price"`
	FinalPrice     int       `json:"finalPrice"`
	SalePrice      int       `json:"salePrice"`
	Comment        string    `json:"comment"`
	WarehouseID    int64     `json:"warehouseId"`
	OfficeID       int       `json:"officeId"`
	Address        *FBSAddress `json:"address,omitempty"`
	Offices        []string  `json:"offices"`
	Skus           []string  `json:"skus"`
	RequiredMeta   []string  `json:"requiredMeta"`
	OptionalMeta   []string  `json:"optionalMeta"`
	ColorCode      string    `json:"colorCode"`
	Rid            string    `json:"rid"`
	Ddate          string    `json:"ddate"`
	SellerDate     string    `json:"sellerDate"`
	ScanPrice      *int      `json:"scanPrice"`
	CurrencyCode   int       `json:"currencyCode"`
	ConvertedPrice int       `json:"convertedPrice"`
	CargoType      int       `json:"cargoType"`
	CrossBorderType int      `json:"crossBorderType"`
	IsZeroOrder    bool      `json:"isZeroOrder"`
}

// FBSAddress — адрес доставки из сборочного задания.
type FBSAddress struct {
	FullAddress string   `json:"fullAddress"`
	Longitude   float64  `json:"longitude"`
	Latitude    float64  `json:"latitude"`
}

// GetOrdersStatusRequest — тело POST /api/v3/orders/status.
type GetOrdersStatusRequest struct {
	Orders []int64 `json:"orders"`
}

// GetOrdersStatusResponse — ответ POST /api/v3/orders/status.
type GetOrdersStatusResponse struct {
	Orders []FBSOrderStatus `json:"orders"`
}

// FBSOrderStatus — статус сборочного задания.
type FBSOrderStatus struct {
	ID             int64  `json:"id"`
	SupplierStatus string `json:"supplierStatus"`
	WBStatus       string `json:"wbStatus"`
}
