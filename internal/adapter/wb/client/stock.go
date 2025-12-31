package client

import (
	"context"
	"fmt"
	"net/http"
	"salehelp/internal/adapter/wb/dto"
)

const maxStocksIDsPerRequest = 1000

func (c *Client) GetStocks(ctx context.Context, apiKey string, warehouseID int64, req dto.GetStocksRequest) (*dto.GetStocksResponse, error) {
	var result dto.GetStocksResponse

	url := fmt.Sprintf("%s/api/v3/stocks/%d", c.marketplaceURL, warehouseID)

	if err := c.doJSON(ctx, http.MethodPost, url, apiKey, RateMarketplace, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Основной “правильный” метод: получаем остатки по chrtIds (size IDs), режем по 1000.
func (c *Client) GetStocksByChrtIDs(ctx context.Context, apiKey string, warehouseID int64, chrtIDs []int64) ([]dto.StockItem, error) {
	if len(chrtIDs) == 0 {
		return nil, nil
	}

	parts := chunkInt64(chrtIDs, maxStocksIDsPerRequest)

	// out capacity — максимум как len(chrtIDs), но фактически WB может вернуть меньше (если чего-то нет на складе)
	out := make([]dto.StockItem, 0, len(chrtIDs))

	for _, part := range parts {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		resp, err := c.GetStocks(ctx, apiKey, warehouseID, dto.GetStocksRequest{
			ChrtIDs: part,
		})
		if err != nil {
			return nil, err
		}

		out = append(out, resp.Stocks...)
	}

	return out, nil
}

// Если удобнее работать через map для быстрого join-а:
func (c *Client) GetStocksByChrtIDsMap(ctx context.Context, apiKey string, warehouseID int64, chrtIDs []int64) (map[int64]dto.StockItem, error) {
	items, err := c.GetStocksByChrtIDs(ctx, apiKey, warehouseID, chrtIDs)
	if err != nil {
		return nil, err
	}

	out := make(map[int64]dto.StockItem, len(items))
	for _, it := range items {
		out[it.ChrtID] = it
	}

	return out, nil
}
