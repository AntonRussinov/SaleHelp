package model

import (
	"salehelp/domain"
	"time"

	"gorm.io/gorm"
)

// Offer представляет основную таблицу offers.
type Offer struct {
	ID                      int            `gorm:"primaryKey;column:id"`
	OfferExternalID         string         `gorm:"column:offer_external_id;type:varchar(100)"`
	SKU                     string         `gorm:"column:sku;type:varchar(100)"`
	ProductID               int            `gorm:"column:product_id;not null"`
	MerchantUID             string         `gorm:"column:merchant_uid;type:varchar(100)"`
	ShopLink                string         `gorm:"column:shop_link;type:varchar(255)"`
	FileID                  string         `gorm:"column:file_id;type:varchar(50)"`
	OfferSource             string         `gorm:"column:offer_source;type:varchar(50)"`
	ReindexReason           string         `gorm:"column:reindex_reason;type:varchar(100)"`
	MatchesLastBatch        bool           `gorm:"column:matches_last_batch"`
	SkipDistinctCheck       bool           `gorm:"column:skip_distinct_check"`
	Priority                int            `gorm:"column:priority"`
	Price                   float64        `gorm:"column:price;type:numeric(10,2);not null"`
	OldPrice                float64        `gorm:"column:old_price;type:numeric(10,2)"`
	Available               bool           `gorm:"column:available;not null"`
	MinPrice                float64        `gorm:"column:min_price;type:numeric(10,2)"`
	MaxPrice                float64        `gorm:"column:max_price;type:numeric(10,2)"`
	ArchivedAt              *time.Time     `gorm:"column:archived_at"`
	Modification            bool           `gorm:"column:modification"`
	ColourCode              string         `gorm:"column:colour_code;type:varchar(50)"`
	EditionDeadline         *time.Time     `gorm:"column:edition_deadline"`
	AnyPickup               bool           `gorm:"column:any_pickup"`
	AnyKaspiDelivery        bool           `gorm:"column:any_kaspi_delivery"`
	AnyKaspiDeliveryLocal   bool           `gorm:"column:any_kaspi_delivery_local"`
	AnyKaspiDeliveryExpress bool           `gorm:"column:any_kaspi_delivery_express"`
	AnyMerchantDelivery     bool           `gorm:"column:any_merchant_delivery"`
	BrandCode               string         `gorm:"column:brand_code;type:varchar(50)"`
	SelectedCity            string         `gorm:"column:selected_city;type:varchar(50)"`
	OfferError              string         `gorm:"column:offer_error;type:text"`
	CreatedAt               time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt               time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt               gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Product    *Product    `gorm:"foreignKey:ProductID;references:ID"`
	Merchant   *Merchant   `gorm:"foreignKey:MerchantUID;references:MerchantUID"`
	Images     []Image     `gorm:"foreignKey:OfferID;references:ID"`
	CityPrices []CityPrice `gorm:"foreignKey:OfferID;references:ID"`
	Stocks     []Stock     `gorm:"foreignKey:OfferID,MerchantUID;references:OfferID,MerchantUID"` //change foreign key to OfferID plus merchant_uid
}

func OfferColumns() []string {
	return []string{
		"id",
		"offer_external_id",
		"sku",
		"product_id",
		"merchant_uid",
		"shop_link",
		"file_id",
		"offer_source",
		"reindex_reason",
		"matches_last_batch",
		"skip_distinct_check",
		"priority",
		"price",
		"old_price",
		"available",
		"min_price",
		"max_price",
		"archived_at",
		"modification",
		"colour_code",
		"edition_deadline",
		"any_pickup",
		"any_kaspi_delivery",
		"any_kaspi_delivery_local",
		"any_kaspi_delivery_express",
		"any_merchant_delivery",
		"brand_code",
		"selected_city",
		"offer_error",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (o *Offer) ToDomain() *domain.Offer {
	if o == nil {
		return nil
	}
	return &domain.Offer{
		ID:                      o.ID,
		OfferExternalID:         o.OfferExternalID,
		SKU:                     o.SKU,
		ProductID:               o.ProductID,
		MerchantUID:             o.MerchantUID,
		ShopLink:                o.ShopLink,
		FileID:                  o.FileID,
		OfferSource:             o.OfferSource,
		ReindexReason:           o.ReindexReason,
		MatchesLastBatch:        o.MatchesLastBatch,
		SkipDistinctCheck:       o.SkipDistinctCheck,
		Priority:                o.Priority,
		Price:                   o.Price,
		OldPrice:                o.OldPrice,
		Available:               o.Available,
		MinPrice:                o.MinPrice,
		MaxPrice:                o.MaxPrice,
		ArchivedAt:              o.ArchivedAt,
		Modification:            o.Modification,
		ColourCode:              o.ColourCode,
		EditionDeadline:         o.EditionDeadline,
		AnyPickup:               o.AnyPickup,
		AnyKaspiDelivery:        o.AnyKaspiDelivery,
		AnyKaspiDeliveryLocal:   o.AnyKaspiDeliveryLocal,
		AnyKaspiDeliveryExpress: o.AnyKaspiDeliveryExpress,
		AnyMerchantDelivery:     o.AnyMerchantDelivery,
		BrandCode:               o.BrandCode,
		SelectedCity:            o.SelectedCity,
		OfferError:              o.OfferError,
		// Преобразование вложенных сущностей:
		Images:     toDomainImages(o.Images),
		CityPrices: toDomainCityPrices(o.CityPrices),
		Stocks:     toDomainStocks(o.Stocks),
	}
}

func toDomainImages(images []Image) []domain.Image {
	res := make([]domain.Image, len(images))
	for i, img := range images {
		res[i] = *img.ToDomain()
	}
	return res
}

func toDomainCityPrices(prices []CityPrice) []domain.CityPrice {
	res := make([]domain.CityPrice, len(prices))
	for i, cp := range prices {
		res[i] = *cp.ToDomain()
	}
	return res
}

func toDomainStocks(stocks []Stock) []domain.Stock {
	res := make([]domain.Stock, len(stocks))
	for i, s := range stocks {
		res[i] = *s.ToDomain()
	}
	return res
}

func (*Offer) TableName() string {
	return "offers"
}
