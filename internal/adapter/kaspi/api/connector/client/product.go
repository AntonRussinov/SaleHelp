package client

import (
	"context"
	"fmt"
	"salehelp/internal/adapter/kaspi/api/models"
)

func (a *Client) GetProductByOrderEntryId(ctx context.Context, orderEntryId string) (*models.ProductResponse, error) {
	var response models.ProductResponse
	if err := a.Do(ctx, GET, fmt.Sprintf("orderentries/%s/product", orderEntryId), nil, "authToken", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Client) GetProductByMasterProductId(ctx context.Context, masterProductId string) (*models.ProductResponse, error) {
	var response models.ProductResponse
	if err := a.Do(ctx, GET, fmt.Sprintf("masterproducts/%s/merchantProduct", masterProductId), nil, "authToken", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Client) GetProduct(ctx context.Context, productId string) (*models.ProductResponse, error) {
	var response models.ProductResponse
	if err := a.Do(ctx, GET, fmt.Sprintf("orderentries/%s/product", productId), nil, "authToken", &response); err != nil {
		return nil, err
	}
	return &response, nil
}
