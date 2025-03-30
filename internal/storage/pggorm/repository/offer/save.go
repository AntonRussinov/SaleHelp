package offer

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

func (repo *Repository) Save(ctx context.Context, offer *domain.Offer) error {
	const fn = "OfferRepository.Save"

	entity := model.Offer{
		ID:                      offer.ID,
		OfferExternalID:         offer.OfferExternalID,
		SKU:                     offer.SKU,
		ProductID:               offer.ProductID,
		MerchantUID:             offer.MerchantUID,
		ShopLink:                offer.ShopLink,
		FileID:                  offer.FileID,
		OfferSource:             offer.OfferSource,
		ReindexReason:           offer.ReindexReason,
		MatchesLastBatch:        offer.MatchesLastBatch,
		SkipDistinctCheck:       offer.SkipDistinctCheck,
		Priority:                offer.Priority,
		Price:                   offer.Price,
		OldPrice:                offer.OldPrice,
		Available:               offer.Available,
		MinPrice:                offer.MinPrice,
		MaxPrice:                offer.MaxPrice,
		ArchivedAt:              offer.ArchivedAt,
		Modification:            offer.Modification,
		ColourCode:              offer.ColourCode,
		EditionDeadline:         offer.EditionDeadline,
		AnyPickup:               offer.AnyPickup,
		AnyKaspiDelivery:        offer.AnyKaspiDelivery,
		AnyKaspiDeliveryLocal:   offer.AnyKaspiDeliveryLocal,
		AnyKaspiDeliveryExpress: offer.AnyKaspiDeliveryExpress,
		AnyMerchantDelivery:     offer.AnyMerchantDelivery,
		BrandCode:               offer.BrandCode,
		SelectedCity:            offer.SelectedCity,
		OfferError:              offer.OfferError,
		// add nested entities
		Product:    ToModeProduct(offer.Product),
		Merchant:   ToModeMerchant(offer.Merchant),
		Images:     ToModelImages(offer.Images),
		CityPrices: ToModeCityPrice(offer.CityPrices),
		Stocks:     ToModeStock(offer.Stocks),
	}

	err := repo.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&entity).Error
	if err != nil {
		return repo.HandleError(fn, err)
	}

	return nil
}

func ToModelImages(request []domain.Image) []model.Image {
	if request == nil {
		return nil
	}

	res := make([]model.Image, len(request))
	for i, img := range request {
		res[i] = model.Image{
			ID:        img.ID,
			OfferID:   img.OfferID,
			ImagePath: img.ImagePath,
		}
	}

	return nil
}

func ToModeProduct(request *domain.Product) *model.Product {
	if request == nil {
		return nil
	}

	return &model.Product{
		ID:               request.ID,
		MasterSKU:        request.MasterSKU,
		Model:            request.Model,            // модель товара
		Brand:            request.Brand,            // бренд
		MasterTitle:      request.MasterTitle,      // основное название товара
		Title:            request.Title,            // дополнительное название или описание
		VerticalCategory: request.VerticalCategory, // вертикальная категория
		MasterCategory:   request.MasterCategory,   // основная категория
	}
}
func ToModeMerchant(request *domain.Merchant) *model.Merchant {
	if request == nil {
		return nil
	}

	return &model.Merchant{
		CompanyName: request.CompanyName,
		MerchantUID: request.MerchantUID,
		IsEnabled:   request.IsEnabled,
	}
}
func ToModeAvailabilities(request []domain.Availability) []model.Availability {
	if request == nil {
		return nil
	}

	res := make([]model.Availability, len(request))
	for i, a := range request {
		res[i] = model.Availability{
			ID:           a.ID,
			MerchantUID:  a.MerchantUID,
			OfferID:      a.OfferID,
			Available:    a.Available,
			StoreID:      a.StoreID,
			PreOrder:     a.PreOrder,
			StockCount:   a.StockCount,
			StockEnabled: a.StockEnabled,
		}
	}

	return res
}
func ToModeCityPrice(request []domain.CityPrice) []model.CityPrice {
	if request == nil {
		return nil
	}

	res := make([]model.CityPrice, len(request))
	for i, a := range request {
		res[i] = model.CityPrice{
			ID:          a.ID,
			OfferID:     a.OfferID,
			CityID:      a.CityID,
			Price:       a.Price,
			OldPrice:    a.OldPrice,
			MerchantUID: a.MerchantUID,
		}
	}

	return res
}
func ToModeStock(request []domain.Stock) []model.Stock {
	if request == nil {
		return nil
	}

	res := make([]model.Stock, len(request))
	for i, a := range request {
		res[i] = model.Stock{
			OfferID:     a.OfferID,
			MerchantUID: a.MerchantUID,
			StoreID:     a.StoreID,
			CityID:      a.CityID,
			Available:   a.Available,
			PreOrder:    a.PreOrder,
			StockCount:  a.StockCount,
			Enabled:     a.Enabled,
		}
	}

	return res
}
