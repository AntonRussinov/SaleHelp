package dto

type GetPricesRequest struct {
	NmList []int64 `json:"nmList"`
}

type GetPricesResponse struct {
	Data struct {
		ListGoods []PriceGood `json:"listGoods"`
	} `json:"data"`
	Error     bool   `json:"error"`
	ErrorText string `json:"errorText"`
}

type PriceGood struct {
	NmID       int64       `json:"nmID"`
	VendorCode string      `json:"vendorCode"`
	Sizes      []PriceSize `json:"sizes"`
	// currencyIsoCode4217
	CurrencyISOCode4217 string  `json:"currencyIsoCode4217"`
	Discount            float64 `json:"discount"`
	ClubDiscount        float64 `json:"clubDiscount"`
	EditableSizePrice   bool    `json:"editableSizePrice"`
	IsBadTurnover       bool    `json:"isBadTurnover"`
}

type PriceSize struct {
	ChrtID              int64   `json:"chrtID"`
	SizeID              int64   `json:"sizeID"`
	Price               float64 `json:"price"`
	DiscountedPrice     float64 `json:"discountedPrice"`
	ClubDiscountedPrice float64 `json:"clubDiscountedPrice"`
	TechSizeName        string  `json:"techSizeName"`
}

// UploadTaskRequest — тело POST /api/v2/upload/task (цены и скидки по товару, editable_size_price = false).
type UploadTaskRequest struct {
	Data []UploadTaskGood `json:"data"`
}

type UploadTaskGood struct {
	NmID     int64   `json:"nmID"`
	Price    float64 `json:"price"`
	Discount int     `json:"discount"`
}

// UploadTaskResponse — ответ upload/task и upload/task/size.
type UploadTaskResponse struct {
	Data struct {
		ID            int64 `json:"id"`
		AlreadyExists bool  `json:"alreadyExists"`
	} `json:"data"`
	Error     bool   `json:"error"`
	ErrorText string `json:"errorText"`
}

// UploadTaskSizeRequest — тело POST /api/v2/upload/task/size (цены по размерам, editable_size_price = true).
type UploadTaskSizeRequest struct {
	Data []UploadTaskSizeGood `json:"data"`
}

type UploadTaskSizeGood struct {
	NmID   int64   `json:"nmID"`
	SizeID int64   `json:"sizeID"`
	Price  float64 `json:"price"`
}
