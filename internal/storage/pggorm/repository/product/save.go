package product

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

func (repo *Repository) Save(ctx context.Context, product *domain.Product) error {
	const fn = "ProductRepository.Save"

	entity := model.Product{
		ID:               product.ID,
		MasterSKU:        product.MasterSKU,
		Model:            product.Model,
		Brand:            product.Brand,
		MasterTitle:      product.MasterTitle,
		Title:            product.Title,
		VerticalCategory: product.VerticalCategory,
		MasterCategory:   product.MasterCategory,
	}

	err := repo.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&entity).Error
	if err != nil {
		return repo.HandleError(fn, err)
	}

	return nil
}
