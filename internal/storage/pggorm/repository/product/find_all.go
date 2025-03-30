package product

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/params/filter"
	"salehelp/internal/storage/pggorm/gormscopes"
	"salehelp/internal/storage/pggorm/model"
)

func (repo *Repository) FindAllByMerchID(ctx context.Context, merchID string) ([]*domain.Product, error) {
	const fn = "ProductRepository.FindAll"

	var entities []model.Product
	//var totalCount int64

	// if err := repo.DB.Model(&model.Offer{}).Joins("User").
	// 	Scopes(gormscopes.Filter(ctx, filterMap()),
	// 		gormscopes.QueryString(ctx, queryString())).
	// 	Count(&totalCount).Error; err != nil {
	// 	return nil, 0, repo.HandleError(fn, err)
	// }

	err := repo.DB.Model(&model.Product{}).
		Scopes(gormscopes.Filter(ctx, filterMap()),
			gormscopes.QueryString(ctx, queryString()),
			gormscopes.Pagination(ctx),
			gormscopes.Order(ctx, orderMap(), "id")).
		Where("products.merchant_uid = ?", merchID).
		Find(&entities).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	products := make([]*domain.Product, len(entities))

	for i, entity := range entities {
		products[i] = entity.ToDomain()
	}

	return products, nil
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
				FROM product_requests lr
			   WHERE lr.user_id = products.id
				 AND lr.status = 'Pending'))`,
			ValueType: filter.ValueTypeBool,
		},
	}
}
