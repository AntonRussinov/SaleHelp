package stock

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/gormscopes"
	"salehelp/internal/storage/pggorm/model"
)

func (repo *Repository) FindAllByMerchID(ctx context.Context, merchID string) ([]*domain.Stock, error) {
	const fn = "StockRepository.FindAll"

	var entities []model.Stock

	err := repo.DB.Model(&model.Stock{}).
		Scopes( // gormscopes.Filter(ctx, filterMap()),
			// gormscopes.QueryString(ctx, queryString()),
			gormscopes.Pagination(ctx),
			gormscopes.Order(ctx, orderMap(), "id")).
		Where("stocks.merchant_uid = ?", merchID).
		Find(&entities).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	stocks := make([]*domain.Stock, len(entities))

	for i, entity := range entities {
		stocks[i] = entity.ToDomain()
	}

	return stocks, nil
}

func (repo *Repository) FindAllByOfferID(ctx context.Context, offerID string) ([]*domain.Stock, error) {
	const fn = "StockRepository.FindAll"

	var entities []model.Stock

	err := repo.DB.Model(&model.Stock{}).
		Scopes( // gormscopes.Filter(ctx, filterMap()),
			// gormscopes.QueryString(ctx, queryString()),
			gormscopes.Pagination(ctx),
			gormscopes.Order(ctx, orderMap(), "id")).
		Where("stocks.offer_id = ?", offerID).
		Find(&entities).Error
	if err != nil {
		return nil, repo.HandleError(fn, err)
	}

	stocks := make([]*domain.Stock, len(entities))

	for i, entity := range entities {
		stocks[i] = entity.ToDomain()
	}

	return stocks, nil
}

// func queryString() string {
// 	return `("User"."name" ILIKE '%' || @queryString || '%'
// 	OR "User"."email" ILIKE '%' || @queryString || '%')`
// }

func orderMap() map[string]string {
	return map[string]string{
		"stock": `"stock_value"`,
	}
}

// func filterMap() gormscopes.FilterMap {
// 	return gormscopes.FilterMap{
// 		"hasRequest": {
// 			FieldName: `(EXISTS(SELECT lr.id
// 				FROM stock_requests lr
// 			   WHERE lr.user_id = stocks.id
// 				 AND lr.status = 'Pending'))`,
// 			ValueType: filter.ValueTypeBool,
// 		},
// 	}
// }
