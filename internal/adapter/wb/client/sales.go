package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"salehelp/internal/adapter/wb/dto"
)

// GetSales возвращает отчёт продаж/возвратов WB (GET /api/v1/supplier/sales).
// dateFrom — дата последнего изменения; для пагинации передавать lastChangeDate последней строки предыдущего ответа.
// reportLoc — пояс магазина для dateFrom (WB отдаёт даты в поясе продавца); nil = c.reportLocation из конфига.
// flag: 0 (по умолчанию) — данные с lastChangeDate >= dateFrom; 1 — все записи за дату dateFrom.
// Лимит WB: 1 запрос в минуту на аккаунт.
func (c *Client) GetSales(ctx context.Context, apiKey string, dateFrom time.Time, flag int, reportLoc *time.Location) (dto.GetSalesResponse, error) {
	u := c.statisticsURL + "/api/v1/supplier/sales"
	params := url.Values{}
	loc := reportLoc
	if loc == nil {
		loc = c.reportLocation
	}
	if loc == nil {
		loc = time.UTC
	}
	params.Set("dateFrom", dateFrom.In(loc).Format(time.RFC3339))
	if flag != 0 {
		params.Set("flag", "1")
	}
	u += "?" + params.Encode()

	var result dto.GetSalesResponse
	if err := c.doJSON(ctx, http.MethodGet, u, apiKey, RateStatistics, nil, &result); err != nil {
		return nil, fmt.Errorf("wb get sales: %w", err)
	}
	return result, nil
}
