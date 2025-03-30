package image

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/params/filter"
	"salehelp/internal/storage/pggorm/gormscopes"
	"salehelp/internal/storage/pggorm/model"
)

func (repo *Repository) FindAllByMerchID(ctx context.Context, offerID string) ([]*domain.Image, error) {
	const fn = "ImageRepository.FindAll"

	var entities []model.Image

	err := repo.DB.Model(&model.Image{}).
		Scopes(gormscopes.Filter(ctx, filterMap()),
			gormscopes.QueryString(ctx, queryString()),
			gormscopes.Pagination(ctx),
			gormscopes.Order(ctx, orderMap(), "id")).
		Where("images.offer_id = ?", offerID).
		Find(&entities).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	images := make([]*domain.Image, len(entities))

	for i, entity := range entities {
		images[i] = entity.ToDomain()
	}

	return images, nil
}

func queryString() string {
	return `("User"."name" ILIKE '%' || @queryString || '%'
	OR "User"."email" ILIKE '%' || @queryString || '%')`
}

func orderMap() map[string]string {
	return map[string]string{
		"name": `"User"."name"`,
	}
}

func filterMap() gormscopes.FilterMap {
	return gormscopes.FilterMap{
		"hasRequest": {
			FieldName: `(EXISTS(SELECT lr.id
				FROM image_requests lr
			   WHERE lr.user_id = images.id
				 AND lr.status = 'Pending'))`,
			ValueType: filter.ValueTypeBool,
		},
	}
}
