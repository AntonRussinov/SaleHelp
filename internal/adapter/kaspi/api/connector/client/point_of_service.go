package client

import (
	"context"
	"fmt"
	"salehelp/internal/adapter/kaspi/api/models"
)

func (a *Client) GetPointOfServiceByOrderEntryId(ctx context.Context, orderEntryId string) (*models.PointOfServiceResponse, error) {
	var response models.PointOfServiceResponse
	if err := a.Do(ctx, GET, fmt.Sprintf("orderentries/%s/deliveryPointOfService", orderEntryId), nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Client) GetPointOfService(ctx context.Context, pointOfServiceId string) (*models.PointOfServiceResponse, error) {
	var response models.PointOfServiceResponse
	if err := a.Do(ctx, GET, fmt.Sprintf("pointofservices/%s", pointOfServiceId), nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
