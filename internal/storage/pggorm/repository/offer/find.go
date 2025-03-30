package offer

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (repo *Repository) FindByProductID(ctx context.Context, id string) (*domain.Offer, error) {
	const fn = "OfferRepository.FindByMerchID"

	var entity model.Offer

	err := repo.DB.WithContext(ctx).First(&entity, "offers.product_id = ?", id).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	cityPrice := entity.ToDomain()

	return cityPrice, nil
}

func (repo *Repository) FindByMerchantID(ctx context.Context, id string) (*domain.Offer, error) {
	const fn = "OfferRepository.FindByMerchantID"

	var entity model.Offer

	err := repo.DB.WithContext(ctx).First(&entity, "offers.merchant_uid = ?", id).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	cityPrice := entity.ToDomain()

	return cityPrice, nil
}
