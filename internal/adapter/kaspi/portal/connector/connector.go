package connector

import (
	"context"
	"salehelp/internal/adapter/kaspi/portal/connector/client"
	"salehelp/internal/adapter/kaspi/portal/portal_models"
)

type PortalAPI interface {
	GetAllProducts(ctx context.Context, merchantID string) (*portal_models.Products, error)
}

func New(MerchantRepo client.MerchantRepo) PortalAPI {
	return client.New(MerchantRepo)
}
