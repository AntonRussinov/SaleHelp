package cityprice

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (repo *Repository) FindByMerchID(ctx context.Context, id string) (*domain.CityPrice, error) {
	const fn = "CityPriceRepository.FindByMerchID"

	var entity model.CityPrice

	err := repo.DB.WithContext(ctx).First(&entity, "city_prices.merchant_uid = ?", id).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	cityPrice := entity.ToDomain()

	return cityPrice, nil
}

func (repo *Repository) FindByOfferID(ctx context.Context, id string) (*domain.CityPrice, error) {
	const fn = "CityPriceRepository.FindByOfferID"

	var entity model.CityPrice

	err := repo.DB.WithContext(ctx).First(&entity, "city_prices.offer_id = ?", id).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	cityPrice := entity.ToDomain()

	return cityPrice, nil
}
