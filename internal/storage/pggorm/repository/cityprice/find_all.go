package cityprice

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/gormscopes"
	"salehelp/internal/storage/pggorm/model"
)

func (repo *Repository) FindAllByMerchID(ctx context.Context, merchID string) ([]*domain.CityPrice, error) {
	const fn = "CityPriceRepository.FindAll"

	var entities []model.CityPrice

	err := repo.DB.Model(&model.CityPrice{}).
		Scopes( // gormscopes.Filter(ctx, filterMap()),
			// gormscopes.QueryString(ctx, queryString()),
			gormscopes.Pagination(ctx),
			gormscopes.Order(ctx, orderMap(), "id")).
		Where("city_prices.merchant_uid = ?", merchID).
		Find(&entities).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	cityPrices := make([]*domain.CityPrice, len(entities))

	for i, entity := range entities {
		cityPrices[i] = entity.ToDomain()
	}

	return cityPrices, nil
}

func (repo *Repository) FindAllByOfferID(ctx context.Context, offerID string) ([]*domain.CityPrice, error) {
	const fn = "CityPriceRepository.FindAll"

	var entities []model.CityPrice

	err := repo.DB.Model(&model.CityPrice{}).
		Scopes( // gormscopes.Filter(ctx, filterMap()),
			// gormscopes.QueryString(ctx, queryString()),
			gormscopes.Pagination(ctx),
			gormscopes.Order(ctx, orderMap(), "id")).
		Where("city_prices.offer_id = ?", offerID).
		Find(&entities).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	cityPrices := make([]*domain.CityPrice, len(entities))

	for i, entity := range entities {
		cityPrices[i] = entity.ToDomain()
	}

	return cityPrices, nil
}

// func queryString() string {
// 	return `("User"."name" ILIKE '%' || @queryString || '%'
// 	OR "User"."email" ILIKE '%' || @queryString || '%')`
// }

func orderMap() map[string]string {
	return map[string]string{
		"cityPrice": `"cityPrice_value"`,
	}
}

// func filterMap() gormscopes.FilterMap {
// 	return gormscopes.FilterMap{
// 		"hasRequest": {
// 			FieldName: `(EXISTS(SELECT lr.id
// 				FROM cityPrice_requests lr
// 			   WHERE lr.user_id = cityPrices.id
// 				 AND lr.status = 'Pending'))`,
// 			ValueType: filter.ValueTypeBool,
// 		},
// 	}
// }
