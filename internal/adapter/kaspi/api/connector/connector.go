package connector

import (
	"context"
	"salehelp/internal/adapter/kaspi/api/connector/client"
	"salehelp/internal/adapter/kaspi/api/models"
)

type API interface {
	GetCities(ctx context.Context) (*models.CitiesResponse, error)
	GetCityByPointOfServiceId(ctx context.Context, pointOfServiceId string) (*models.CityResponse, error)

	GetOrders(ctx context.Context, req models.GetOrdersRequest) (*models.OrdersResponse, error)
	GetOrderByCode(ctx context.Context, code string) (*models.OrderByCodeResponse, error)

	GetOrderEntries(ctx context.Context, orderId string) (*models.OrderEntriesResponse, error)
	GetOrderEntry(ctx context.Context, entryId string) (*models.OrderEntryResponse, error)

	GetPointOfServiceByOrderEntryId(ctx context.Context, orderEntryId string) (*models.PointOfServiceResponse, error)
	GetPointOfService(ctx context.Context, pointOfServiceId string) (*models.PointOfServiceResponse, error)

	GetProductByOrderEntryId(ctx context.Context, orderEntryId string) (*models.ProductResponse, error)
	GetProductByMasterProductId(ctx context.Context, masterProductId string) (*models.ProductResponse, error)
	GetProduct(ctx context.Context, productId string) (*models.ProductResponse, error)

	GetReviews(ctx context.Context, req models.GetReviewsRequest) (*models.ReviewsResponse, error)
	GetReviewByOrderId(ctx context.Context, orderId string) (*models.ReviewResponse, error)
}

func New(authToken string) API {
	return client.New(authToken)
}
