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
