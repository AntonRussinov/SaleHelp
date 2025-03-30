package client

import (
	"context"
	"net/url"
	"salehelp/internal/adapter/kaspi/api/models"
)

func (a *Client) GetOrders(ctx context.Context, req models.GetOrdersRequest) (*models.OrdersResponse, error) {
	reqParams, err := req.ToUrlValues()
	if err != nil {
		return nil, err
	}

	var response models.OrdersResponse
	if err = a.Do(ctx, GET, "orders", reqParams, "authToken", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Client) GetOrderByCode(ctx context.Context, code string) (*models.OrderByCodeResponse, error) {
	var response models.OrderByCodeResponse
	if err := a.Do(ctx, GET, "orders", url.Values{"filter[orders][code]": []string{code}}, "authToken", &response); err != nil {
		return nil, err
	}
	return &response, nil
}
