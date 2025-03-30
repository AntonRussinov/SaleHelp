package stock

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

func (repo *Repository) Save(ctx context.Context, stock *domain.Stock) error {
	const fn = "StockRepository.Save"

	entity := model.Stock{
		MerchantUID: stock.MerchantUID,
		OfferID:     stock.OfferID,
		CityID:      stock.CityID,
		StoreID:     stock.StoreID,
		StockCount:  stock.StockCount,
		Available:   stock.Available,
		PreOrder:    stock.PreOrder,
		Enabled:     stock.Enabled,
	}

	err := repo.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&entity).Error
	if err != nil {
		return repo.HandleError(fn, err)
	}

	return nil
}
