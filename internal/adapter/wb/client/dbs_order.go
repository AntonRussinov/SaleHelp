package client

import (
	"context"
	"net/http"
	"salehelp/internal/adapter/wb/dto"
)

// GetNewDBSOrders возвращает список новых сборочных заданий DBS (GET /api/v3/dbs/orders/new).
// Лимит WB: 300 запросов/мин, интервал 200 мс (RateMarketplace).
func (c *Client) GetNewDBSOrders(ctx context.Context, apiKey string) (*dto.GetNewDBSOrdersResponse, error) {
	var result dto.GetNewDBSOrdersResponse
	url := c.marketplaceURL + "/api/v3/dbs/orders/new"
	if err := c.doJSON(ctx, http.MethodGet, url, apiKey, RateMarketplace, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
