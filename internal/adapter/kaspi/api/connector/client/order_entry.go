package client

import (
	"context"
	"fmt"
	"salehelp/internal/adapter/kaspi/api/models"
)

func (a *Client) GetOrderEntries(ctx context.Context, orderId string) (*models.OrderEntriesResponse, error) {
	var response models.OrderEntriesResponse
	if err := a.Do(ctx, GET, fmt.Sprintf("orders/%s/entries", orderId), nil, "authToken", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Client) GetOrderEntry(ctx context.Context, entryId string) (*models.OrderEntryResponse, error) {
	var response models.OrderEntryResponse
	if err := a.Do(ctx, GET, fmt.Sprintf("orderentries/%s", entryId), nil, "authToken", &response); err != nil {
		return nil, err
	}
	return &response, nil
}
