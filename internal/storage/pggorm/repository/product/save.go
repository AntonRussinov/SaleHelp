package product

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

func (repo *Repository) Save(ctx context.Context, stock *domain.Product) error {
	const fn = "ProductRepository.Save"

	entity := model.Product{
		ID:               stock.ID,
		MasterSKU:        stock.MasterSKU,
		Model:            stock.Model,
		Brand:            stock.Brand,
		MasterTitle:      stock.MasterTitle,
		Title:            stock.Title,
		VerticalCategory: stock.VerticalCategory,
		MasterCategory:   stock.MasterCategory,
	}

	err := repo.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&entity).Error
	if err != nil {
		return repo.HandleError(fn, err)
	}

	return nil
}
