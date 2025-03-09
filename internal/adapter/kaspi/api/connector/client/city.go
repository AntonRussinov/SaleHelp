package client

import (
	"context"
	"fmt"
	"salehelp/internal/adapter/kaspi/api/models"
)

func (a *Client) GetCities(ctx context.Context) (*models.CitiesResponse, error) {
	var response models.CitiesResponse
	if err := a.Do(ctx, GET, "cities", nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Client) GetCityByPointOfServiceId(ctx context.Context, pointOfServiceId string) (*models.CityResponse, error) {
	var response models.CityResponse
	if err := a.Do(ctx, GET, fmt.Sprintf("pointofservices/%s/city", pointOfServiceId), nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
