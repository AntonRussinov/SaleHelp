package connector

import (
	"context"
	"salehelp/internal/adapter/kaspi/portal/connector/client"
	"salehelp/internal/adapter/kaspi/portal/models"
)

type Portal interface {
	GetAllProducts(ctx context.Context, merchantID string) (*models.Products, error)
}

func New(mcsid, mcsession string) Portal {
	return client.New(mcsid, mcsession)
}
