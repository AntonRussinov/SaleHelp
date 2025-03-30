package merchant

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (repo *Repository) FindByID(ctx context.Context, id string) (*domain.Merchant, error) {
	const fn = "MerchantRepository.FindByID"

	var entity model.Merchant

	err := repo.DB.WithContext(ctx).First(&entity, "merchants.merchant_uid = ?", id).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	merchant := entity.ToDomain()

	return merchant, nil
}
