package client

import (
	"context"
	"net/http"
	"salehelp/internal/adapter/wb/dto"
)

// GetNewOrders возвращает список новых сборочных заданий FBS (GET /api/v3/orders/new).
// Лимит WB (базовый токен): 30/час.
func (c *Client) GetNewOrders(ctx context.Context, apiKey string) (*dto.GetNewOrdersResponse, error) {
	var result dto.GetNewOrdersResponse
	url := c.marketplaceURL + "/api/v3/orders/new"
	if err := c.doJSON(ctx, http.MethodGet, url, apiKey, RateFBSOrders, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOrdersStatus возвращает статусы сборочных заданий по ID (POST /api/v3/orders/status).
// orders — до 1000 ID. Лимит WB (базовый токен): 30/час.
func (c *Client) GetOrdersStatus(ctx context.Context, apiKey string, orderIDs []int64) (*dto.GetOrdersStatusResponse, error) {
	if len(orderIDs) == 0 {
		return &dto.GetOrdersStatusResponse{Orders: nil}, nil
	}
	var result dto.GetOrdersStatusResponse
	url := c.marketplaceURL + "/api/v3/orders/status"
	req := dto.GetOrdersStatusRequest{Orders: orderIDs}
	if err := c.doJSON(ctx, http.MethodPost, url, apiKey, RateFBSOrders, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
