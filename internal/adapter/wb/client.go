package wb

import (
	"context"
	"salehelp/internal/adapter/wb/dto"
)

type Client interface {
	GetProducts(ctx context.Context, apiKey string) ([]dto.Card, error)
	// GetStocks(ctx context.Context, apiKey string) ([]StockDTO, error)
	// GetPrices(ctx context.Context, apiKey string) ([]PriceDTO, error)
}
