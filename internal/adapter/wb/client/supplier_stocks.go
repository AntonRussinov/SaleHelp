package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"salehelp/internal/adapter/wb/dto"
)

// GetSupplierStocks возвращает текущие остатки на складах WB (GET /api/v1/supplier/stocks).
// Включает товары от возвратов, которые WB хранит и продаёт самостоятельно.
// dateFrom — дата начала периода; для получения актуальных остатков передавай дату ~30–90 дней назад.
// Лимит WB: 1 запрос в минуту (RateStatistics).
func (c *Client) GetSupplierStocks(ctx context.Context, apiKey string, dateFrom time.Time) (dto.GetSupplierStocksResponse, error) {
	u := c.statisticsURL + "/api/v1/supplier/stocks"
	params := url.Values{}
	params.Set("dateFrom", dateFrom.Format("2006-01-02"))
	u += "?" + params.Encode()

	var result dto.GetSupplierStocksResponse
	if err := c.doJSON(ctx, http.MethodGet, u, apiKey, RateStatistics, nil, &result); err != nil {
		return nil, fmt.Errorf("wb get supplier stocks: %w", err)
	}
	return result, nil
}
