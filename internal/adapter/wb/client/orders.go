package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"salehelp/internal/adapter/wb/dto"
)

// GetOrders возвращает отчёт заказов WB (GET /api/v1/supplier/orders).
// Связь с Sales по srid. isCancel=true — заказ не забран (отказ).
// dateFrom — дата последнего изменения; для пагинации — lastChangeDate последней строки предыдущего ответа.
// flag: 0 — данные с lastChangeDate >= dateFrom; 1 — все записи за дату dateFrom.
// Лимит WB: 1 запрос в минуту на аккаунт (Statistics).
func (c *Client) GetOrders(ctx context.Context, apiKey string, dateFrom time.Time, flag int) (dto.GetOrdersResponse, error) {
	u := c.statisticsURL + "/api/v1/supplier/orders"
	params := url.Values{}
	params.Set("dateFrom", dateFrom.Format(time.RFC3339))
	if flag != 0 {
		params.Set("flag", "1")
	}
	u += "?" + params.Encode()

	var result dto.GetOrdersResponse
	if err := c.doJSON(ctx, http.MethodGet, u, apiKey, RateStatistics, nil, &result); err != nil {
		return nil, fmt.Errorf("wb get orders: %w", err)
	}
	return result, nil
}
