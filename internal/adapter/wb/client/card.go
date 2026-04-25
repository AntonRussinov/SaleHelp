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

// GetAllCards загружает карточки товаров постранично.
//
// startCursor — стартовая позиция для инкрементального синка (updatedAt + nmID из предыдущего ответа WB).
// Если nil — полная выгрузка с самого начала.
//
// Возвращает список карточек и финальный курсор (для сохранения после синка).
// Если карточек нет (ничего не изменилось с момента cursor), возвращает пустой список и nil cursor.
func (c *Client) GetAllCards(ctx context.Context, apiKey string, baseReq dto.GetCardsRequest, startCursor *dto.ResponseCursor) ([]dto.Card, *dto.ResponseCursor, error) {
	limit := baseReq.Settings.Cursor.Limit
	if limit <= 0 {
		limit = 100
		baseReq.Settings.Cursor.Limit = limit
	}

	// Применяем стартовый курсор для инкрементального синка.
	if startCursor != nil && startCursor.UpdatedAt != "" {
		baseReq.Settings.Cursor.UpdatedAt = startCursor.UpdatedAt
		baseReq.Settings.Cursor.NmID = startCursor.NmID
	}

	const maxPages = 10_000

	var all []dto.Card
	var finalCursor *dto.ResponseCursor
	req := baseReq

	for page := 1; page <= maxPages; page++ {
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}

		resp, err := c.GetCards(ctx, apiKey, req)
		if err != nil {
			return nil, nil, err
		}

		if len(resp.Cards) > 0 {
			all = append(all, resp.Cards...)
			// Сохраняем курсор с каждой страницы — нужен финальный.
			cur := resp.Cursor
			finalCursor = &cur
		}

		if resp.Cursor.Total < limit || len(resp.Cards) == 0 {
			break
		}

		req.Settings.Cursor.UpdatedAt = resp.Cursor.UpdatedAt
		req.Settings.Cursor.NmID = resp.Cursor.NmID
	}

	return all, finalCursor, nil
}
