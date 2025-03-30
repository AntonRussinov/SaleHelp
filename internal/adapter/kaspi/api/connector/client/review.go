package client

import (
	"context"
	"fmt"
	"salehelp/internal/adapter/kaspi/api/models"
)

func (a *Client) GetReviews(ctx context.Context, req models.GetReviewsRequest) (*models.ReviewsResponse, error) {
	reqParams, err := req.ToUrlValues()
	if err != nil {
		return nil, err
	}

	var response models.ReviewsResponse
	if err = a.Do(ctx, GET, "merchantreviews", reqParams, "authToken", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *Client) GetReviewByOrderId(ctx context.Context, orderId string) (*models.ReviewResponse, error) {
	var response models.ReviewResponse
	if err := a.Do(ctx, GET, fmt.Sprintf("merchantreviews/%s", orderId), nil, "authToken", &response); err != nil {
		return nil, err
	}
	return &response, nil
}
