package cityprice

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

func (repo *Repository) Save(ctx context.Context, cityPrice *domain.CityPrice) error {
	const fn = "CityPriceRepository.Save"

	entity := model.CityPrice{
		ID:          cityPrice.ID,
		MerchantUID: cityPrice.MerchantUID,
		OfferID:     cityPrice.OfferID,
		CityID:      cityPrice.CityID,
		Price:       cityPrice.Price,
		OldPrice:    cityPrice.OldPrice,
	}

	err := repo.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&entity).Error
	if err != nil {
		return repo.HandleError(fn, err)
	}

	return nil
}
