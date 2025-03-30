package connector

import (
	"context"
	"salehelp/internal/adapter/kaspi/portal/connector/client"
	"salehelp/internal/adapter/kaspi/portal/models"
)

type PortalAPI interface {
	GetAllProducts(ctx context.Context, merchantID string) (*models.Products, error)
}

func New() PortalAPI {
	return client.New()
}
