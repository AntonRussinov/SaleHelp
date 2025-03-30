package image

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

func (repo *Repository) Save(ctx context.Context, iamge *domain.Image) error {
	const fn = "ImageRepository.Save"

	entity := model.Image{
		ID:        iamge.ID,
		ImagePath: iamge.ImagePath,
		OfferID:   iamge.OfferID,
	}

	err := repo.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&entity).Error
	if err != nil {
		return repo.HandleError(fn, err)
	}

	return nil
}
