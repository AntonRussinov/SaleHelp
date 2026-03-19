package client

import (
	"context"
	"fmt"
	"net/http"
	"salehelp/internal/adapter/wb/dto"
)

const maxPriceNmIDsPerRequest = 1000

func (c *Client) GetPrices(ctx context.Context, apiKey string, req dto.GetPricesRequest) (*dto.GetPricesResponse, error) {
	var result dto.GetPricesResponse

	url := c.priceURL + "/api/v2/list/goods/filter"

	if err := c.doJSON(ctx, http.MethodPost, url, apiKey, RatePrices, req, &result); err != nil {
		return nil, err
	}

	// Важно: WB может вернуть 200, но error=true
	if result.Error {
		if result.ErrorText != "" {
			return nil, fmt.Errorf("wb get prices error: %s", result.ErrorText)
		}
		return nil, fmt.Errorf("wb get prices error: unknown")
	}

	return &result, nil
}

func (c *Client) GetPricesByNmIDs(ctx context.Context, apiKey string, nmIDs []int64) ([]dto.PriceGood, error) {
	if len(nmIDs) == 0 {
		return nil, nil
	}

	parts := chunkInt64(nmIDs, maxPriceNmIDsPerRequest)

	// Прикидываем capacity, чтобы меньше аллокаций
	out := make([]dto.PriceGood, 0, len(nmIDs))

	for _, part := range parts {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		resp, err := c.GetPrices(ctx, apiKey, dto.GetPricesRequest{NmList: part})
		if err != nil {
			return nil, err
		}

		out = append(out, resp.Data.ListGoods...)
	}

	return out, nil
}

func (c *Client) GetPricesByNmIDsMap(ctx context.Context, apiKey string, nmIDs []int64) (map[int64]dto.PriceGood, error) {
	goods, err := c.GetPricesByNmIDs(ctx, apiKey, nmIDs)
	if err != nil {
		return nil, err
	}

	out := make(map[int64]dto.PriceGood, len(goods))
	for _, g := range goods {
		out[g.NmID] = g
	}

	return out, nil
}

// UploadPrices устанавливает цены и скидки для товаров (editable_size_price = false).
// POST /api/v2/upload/task. Максимум 1000 товаров в одном запросе.
func (c *Client) UploadPrices(ctx context.Context, apiKey string, req dto.UploadTaskRequest) (*dto.UploadTaskResponse, error) {
	if len(req.Data) == 0 {
		return &dto.UploadTaskResponse{}, nil
	}
	var result dto.UploadTaskResponse
	url := c.priceURL + "/api/v2/upload/task"
	if err := c.doJSON(ctx, http.MethodPost, url, apiKey, RatePrices, req, &result); err != nil {
		return nil, err
	}
	if result.Error && result.ErrorText != "" {
		return &result, fmt.Errorf("wb upload prices: %s", result.ErrorText)
	}
	return &result, nil
}

// UploadPricesBySize устанавливает цены по размерам (editable_size_price = true).
// POST /api/v2/upload/task/size. Максимум 1000 размеров в одном запросе.
func (c *Client) UploadPricesBySize(ctx context.Context, apiKey string, req dto.UploadTaskSizeRequest) (*dto.UploadTaskResponse, error) {
	if len(req.Data) == 0 {
		return &dto.UploadTaskResponse{}, nil
	}
	var result dto.UploadTaskResponse
	url := c.priceURL + "/api/v2/upload/task/size"
	if err := c.doJSON(ctx, http.MethodPost, url, apiKey, RatePrices, req, &result); err != nil {
		return nil, err
	}
	if result.Error && result.ErrorText != "" {
		return &result, fmt.Errorf("wb upload prices by size: %s", result.ErrorText)
	}
	return &result, nil
}
