package client

import (
	"context"
	"net/http"
	"salehelp/internal/adapter/wb/dto"
)

func (c *Client) GetWBWarehouses(ctx context.Context, apiKey string) ([]dto.WBWarehouse, error) {
	url := c.marketplaceURL + "/api/v3/offices"

	var result []dto.WBWarehouse
	if err := c.doJSON(ctx, http.MethodGet, url, apiKey, RateMarketplace, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetSellerWarehouses(ctx context.Context, apiKey string) ([]dto.SellerWarehouse, error) {
	url := c.marketplaceURL + "/api/v3/warehouses"

	var result []dto.SellerWarehouse
	if err := c.doJSON(ctx, http.MethodGet, url, apiKey, RateMarketplace, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}
