package client

import (
	"context"
	"net/http"
	"salehelp/internal/adapter/wb/dto"
)

// GetNewDBWOrders возвращает список новых сборочных заданий DBW (GET /api/v3/dbw/orders/new).
// Лимит WB (базовый токен): 10/час.
func (c *Client) GetNewDBWOrders(ctx context.Context, apiKey string) (*dto.GetNewDBWOrdersResponse, error) {
	var result dto.GetNewDBWOrdersResponse
	url := c.marketplaceURL + "/api/v3/dbw/orders/new"
	if err := c.doJSON(ctx, http.MethodGet, url, apiKey, RateDBOrders, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
