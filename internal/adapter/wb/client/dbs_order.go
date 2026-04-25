package client

import (
	"context"
	"net/http"
	"salehelp/internal/adapter/wb/dto"
)

// GetNewDBSOrders возвращает список новых сборочных заданий DBS (GET /api/v3/dbs/orders/new).
// Лимит WB (базовый токен): 10/час.
func (c *Client) GetNewDBSOrders(ctx context.Context, apiKey string) (*dto.GetNewDBSOrdersResponse, error) {
	var result dto.GetNewDBSOrdersResponse
	url := c.marketplaceURL + "/api/v3/dbs/orders/new"
	if err := c.doJSON(ctx, http.MethodGet, url, apiKey, RateDBOrders, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
