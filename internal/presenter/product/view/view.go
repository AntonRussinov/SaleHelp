package view

import (
	"salehelp/internal/adapter/kaspi/portal/models"
)

type ProductView struct {
	Data  []OfferView `json:"data"`
	Total int         `json:"total"`
}

// Структура, описывающая один товар/предложение
type OfferView struct {
	MasterSku               string                         `json:"masterSku"`
	Images                  []string                       `json:"images"`
	ShopLink                string                         `json:"shopLink"`
	MasterTitle             string                         `json:"masterTitle"`
	OfferId                 string                         `json:"offerId"`
	FileId                  string                         `json:"fileId"`
	MerchantUid             string                         `json:"merchantUid"`
	Sku                     string                         `json:"sku"`
	Model                   string                         `json:"model"`
	Brand                   *string                        `json:"brand"` // может быть null
	Price                   float64                        `json:"price"`
	Oldprice                float64                        `json:"oldprice"`
	Availabilities          []AvailabilityView             `json:"availabilities"`
	CityPrices              []CityPriceView                `json:"cityPrices"`
	DeliveryOptions         interface{}                    `json:"deliveryOptions"` // может быть null
	Version                 interface{}                    `json:"_version"`        // может быть null
	DistinctVersion         int                            `json:"_distinctVersion"`
	OfferSource             string                         `json:"offerSource"`
	ReindexReason           string                         `json:"reindexReason"`
	MatchesLastBatch        bool                           `json:"matchesLastBatch"`
	SkipDistinctCheck       bool                           `json:"skipDistinctCheck"`
	Priority                int                            `json:"priority"`
	RangePrice              RangePriceView                 `json:"rangePrice"`
	Updates                 []UpdateView                   `json:"updates"`
	ArchivedAt              interface{}                    `json:"archivedAt"` // может быть null
	Modification            bool                           `json:"modification"`
	ColourCode              *string                        `json:"colourCode"`      // может быть null
	Attributes              []interface{}                  `json:"attributes"`      // массив, может быть уточнён по необходимости
	EditionDeadline         interface{}                    `json:"editionDeadline"` // может быть null
	VerticalCategory        string                         `json:"verticalCategory"`
	MasterCategory          string                         `json:"masterCategory"`
	Title                   string                         `json:"title"`
	Available               bool                           `json:"available"`
	MinPrice                float64                        `json:"minPrice"`
	MaxPrice                float64                        `json:"maxPrice"`
	Points                  []string                       `json:"points"`
	UpdatedAt               models.CustomTime              `json:"updatedAt"`
	AnyPickup               bool                           `json:"anyPickup"`
	AnyKaspiDelivery        bool                           `json:"anyKaspiDelivery"`
	AnyKaspiDeliveryLocal   bool                           `json:"anyKaspiDeliveryLocal"`
	AnyKaspiDeliveryExpress bool                           `json:"anyKaspiDeliveryExpress"`
	AnyMerchantDelivery     bool                           `json:"anyMerchantDelivery"`
	BrandCode               *string                        `json:"brandCode"`    // может быть null
	SelectedCity            *string                        `json:"selectedCity"` // может быть null
	AllCityPrices           map[string]CityPriceDetailView `json:"allCityPrices"`
	Stocks                  []StockView                    `json:"stocks"`
	OfferError              interface{}                    `json:"offerError"` // может быть null
}

// Структура для availabilities
type AvailabilityView struct {
	Available    string `json:"available"`
	StoreId      string `json:"storeId"`
	StockCount   int    `json:"stockCount"`
	PreOrder     int    `json:"preOrder"`
	StockEnabled bool   `json:"stockEnabled"`
}

// Структура для cityPrices (массив в основном объекте)
type CityPriceView struct {
	Value    float64 `json:"value"`
	CityId   string  `json:"cityId"`
	Oldprice float64 `json:"oldprice"`
}

// Структура для диапазона цен
type RangePriceView struct {
	MIN float64 `json:"MIN"`
	MAX float64 `json:"MAX"`
}

// Структура для updates
type UpdateView struct {
	OperationType string            `json:"operationType"`
	UploadedBy    string            `json:"uploadedBy"`
	UploadedDate  models.CustomTime `json:"uploadedDate"`
}

// Структура для allCityPrices – цены по городам, представленные в виде мапы
type CityPriceDetailView struct {
	Price    float64  `json:"price"`
	OldPrice *float64 `json:"oldPrice"` // может быть null
}

// Структура для stocks
type StockView struct {
	CityId     string                          `json:"cityId"`
	StockLevel map[string]StockLevelDetailView `json:"stockLevel"`
}

// Детали уровня запасов для каждого storeId в stocks
type StockLevelDetailView struct {
	Value   int  `json:"value"`
	Enabled bool `json:"enabled"`
}
