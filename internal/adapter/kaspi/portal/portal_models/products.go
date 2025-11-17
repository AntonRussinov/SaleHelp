package portal_models

import (
	"fmt"
	"strings"
	"time"
)

type CustomTime time.Time

const customTimeFormat = "2006-01-02T15:04:05.000"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		return nil
	}

	// Сначала пробуем парсить по RFC3339Nano (с временной зоной)
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		*ct = CustomTime(t)
		return nil
	}

	// Попытка парсинга без временной зоны с переменной точностью дробной части
	const layout = "2006-01-02T15:04:05.999999999"
	if t, err := time.Parse(layout, s); err == nil {
		*ct = CustomTime(t)
		return nil
	}

	return fmt.Errorf("unable to parse time: %s", s)
}

func (ct CustomTime) String() string {
	return time.Time(ct).Format(customTimeFormat)
}

// Корневой объект JSON
type Products struct {
	Data  []Offer `json:"data"`
	Total int     `json:"total"`
}

// Структура, описывающая один товар/предложение
type Offer struct {
	MasterSku               string                     `json:"masterSku"`
	Images                  []string                   `json:"images"`
	ShopLink                string                     `json:"shopLink"`
	MasterTitle             string                     `json:"masterTitle"`
	OfferId                 string                     `json:"offerId"`
	FileId                  string                     `json:"fileId"`
	MerchantUid             string                     `json:"merchantUid"`
	Sku                     string                     `json:"sku"`
	Model                   string                     `json:"model"`
	Brand                   *string                    `json:"brand"` // может быть null
	Price                   float64                    `json:"price"`
	Oldprice                float64                    `json:"oldprice"`
	Availabilities          []Availability             `json:"availabilities"`
	CityPrices              []CityPrice                `json:"cityPrices"`
	DeliveryOptions         interface{}                `json:"deliveryOptions"` // может быть null
	Version                 interface{}                `json:"_version"`        // может быть null
	DistinctVersion         int                        `json:"_distinctVersion"`
	OfferSource             string                     `json:"offerSource"`
	ReindexReason           string                     `json:"reindexReason"`
	MatchesLastBatch        bool                       `json:"matchesLastBatch"`
	SkipDistinctCheck       bool                       `json:"skipDistinctCheck"`
	Priority                int                        `json:"priority"`
	RangePrice              RangePrice                 `json:"rangePrice"`
	Updates                 []Update                   `json:"updates"`
	ArchivedAt              interface{}                `json:"archivedAt"` // может быть null
	Modification            bool                       `json:"modification"`
	ColourCode              *string                    `json:"colourCode"`      // может быть null
	Attributes              []interface{}              `json:"attributes"`      // массив, может быть уточнён по необходимости
	EditionDeadline         interface{}                `json:"editionDeadline"` // может быть null
	VerticalCategory        string                     `json:"verticalCategory"`
	MasterCategory          string                     `json:"masterCategory"`
	Title                   string                     `json:"title"`
	Available               bool                       `json:"available"`
	MinPrice                float64                    `json:"minPrice"`
	MaxPrice                float64                    `json:"maxPrice"`
	Points                  []string                   `json:"points"`
	UpdatedAt               CustomTime                 `json:"updatedAt"`
	AnyPickup               bool                       `json:"anyPickup"`
	AnyKaspiDelivery        bool                       `json:"anyKaspiDelivery"`
	AnyKaspiDeliveryLocal   bool                       `json:"anyKaspiDeliveryLocal"`
	AnyKaspiDeliveryExpress bool                       `json:"anyKaspiDeliveryExpress"`
	AnyMerchantDelivery     bool                       `json:"anyMerchantDelivery"`
	BrandCode               *string                    `json:"brandCode"`    // может быть null
	SelectedCity            *string                    `json:"selectedCity"` // может быть null
	AllCityPrices           map[string]CityPriceDetail `json:"allCityPrices"`
	Stocks                  []Stock                    `json:"stocks"`
	OfferError              interface{}                `json:"offerError"` // может быть null
}

// Структура для availabilities
type Availability struct {
	Available    string `json:"available"`
	StoreId      string `json:"storeId"`
	StockCount   int    `json:"stockCount"`
	PreOrder     int    `json:"preOrder"`
	StockEnabled bool   `json:"stockEnabled"`
}

// Структура для cityPrices (массив в основном объекте)
type CityPrice struct {
	Value    float64 `json:"value"`
	CityId   string  `json:"cityId"`
	Oldprice float64 `json:"oldprice"`
}

// Структура для диапазона цен
type RangePrice struct {
	MIN float64 `json:"MIN"`
	MAX float64 `json:"MAX"`
}

// Структура для updates
type Update struct {
	OperationType string     `json:"operationType"`
	UploadedBy    string     `json:"uploadedBy"`
	UploadedDate  CustomTime `json:"uploadedDate"`
}

// Структура для allCityPrices – цены по городам, представленные в виде мапы
type CityPriceDetail struct {
	Price    float64  `json:"price"`
	OldPrice *float64 `json:"oldPrice"` // может быть null
}

// Структура для stocks
type Stock struct {
	CityId     string                      `json:"cityId"`
	StockLevel map[string]StockLevelDetail `json:"stockLevel"`
}

// Детали уровня запасов для каждого storeId в stocks
type StockLevelDetail struct {
	Value   int  `json:"value"`
	Enabled bool `json:"enabled"`
}
