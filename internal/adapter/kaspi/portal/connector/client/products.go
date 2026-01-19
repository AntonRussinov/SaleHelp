package client

import (
	"context"
	"fmt"
	"salehelp/internal/adapter/kaspi/portal/portal_models"
)

const (
	productCount = 10
)

func (a *Client) GetProductsByPage(ctx context.Context, merchantID string, lowStock bool) (*portal_models.Products, error) {
	var response portal_models.Products
	if err := a.Do(ctx, GET, fmt.Sprintf("bff/offer-view/list?m=%v&p=0&l=%v&a=true&lowStock=false", merchantID, productCount), nil, merchantID, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Client) GetAllProducts(ctx context.Context, merchantID string) (*portal_models.Products, error) {
	var response portal_models.Products

	if err := a.Do(ctx, GET, fmt.Sprintf("bff/offer-view/list?m=%v&p=%v&l=%v&a=true&lowStock=false", merchantID, 0, productCount), nil, merchantID, &response); err != nil {
		return nil, err
	}

	// временно, так как сейчас иначе будет авторизация каждый круг
	// if response.Total > productCount {
	// 	for i := 1; i <= response.Total/productCount; i++ {
	// 		var page portal_models.Products
	// 		if err := a.Do(ctx, GET, fmt.Sprintf("bff/offer-view/list?m=%v&p=%v&l=%v&a=true&lowStock=false", merchantID, i, productCount), nil, merchantID, &page); err != nil {
	// 			return nil, err
	// 		}
	// 		response.Data = append(response.Data, page.Data...)
	// 	}
	// }

	return &response, nil
}
