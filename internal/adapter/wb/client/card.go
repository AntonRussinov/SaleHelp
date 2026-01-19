package client

import (
	"context"
	"net/http"
	"salehelp/internal/adapter/wb/dto"
)

func (c *Client) GetCards(ctx context.Context, apiKey string, req dto.GetCardsRequest) (*dto.GetCardsResponse, error) {
	var result dto.GetCardsResponse

	url := c.baseURL + "/content/v2/get/cards/list"

	if err := c.doJSON(ctx, http.MethodPost, url, apiKey, RateContent, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetAllCards(ctx context.Context, apiKey string, baseReq dto.GetCardsRequest) ([]dto.Card, error) {
	limit := baseReq.Settings.Cursor.Limit
	if limit <= 0 {
		limit = 100 // дефолт — безопасный, потом можно вынести в конфиг
		baseReq.Settings.Cursor.Limit = limit
	}

	// Защита от бесконечного цикла (на случай, если WB начнет возвращать странный cursor/total)
	const maxPages = 10_000

	var all []dto.Card
	req := baseReq

	for page := 1; page <= maxPages; page++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		resp, err := c.GetCards(ctx, apiKey, req)
		if err != nil {
			return nil, err
		}

		if len(resp.Cards) > 0 {
			all = append(all, resp.Cards...)
		}

		// Условие окончания по WB для этого метода:
		// cursor.total < cursor.limit  => это последняя страница
		// На практике добавим и страховку len(cards)==0.
		if resp.Cursor.Total < limit || len(resp.Cards) == 0 {
			break
		}

		// Следующий запрос: updatedAt + nmID из ответа
		req.Settings.Cursor.UpdatedAt = resp.Cursor.UpdatedAt
		req.Settings.Cursor.NmID = resp.Cursor.NmID
	}

	return all, nil
}
