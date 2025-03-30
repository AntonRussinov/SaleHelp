package merchant

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

func (repo *Repository) Save(ctx context.Context, merchant *domain.Merchant) error {
	const fn = "MerchantRepository.Save"

	entity := model.Merchant{
		CompanyName: merchant.CompanyName,
		IsEnabled:   merchant.IsEnabled,
		MerchantUID: merchant.MerchantUID,
	}

	err := repo.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&entity).Error
	if err != nil {
		return repo.HandleError(fn, err)
	}

	return nil
}
