package di

import (
	"errors"
	"salehelp/internal/handler/product"
	"salehelp/internal/handler/shop"
)

type Handlers struct {
	Product *product.Handler
	Shop    *shop.Handler
}

var ErrInitHandlerss = errors.New("failed to initialize Handlerss")

func (prov *Provider) InitHandlers() error {
	// authHandlers, err := auth.NewHandlers(prov.cfg.Auth, prov.Services.User)
	// if err != nil {
	// 	return fmt.Errorf("%w: %w", ErrInitHandlerss, err)
	// }

	// prov.Handlerss.Auth = authHandlers

	// temporary
	// prov.Handlers.Product = product.NewHandler(prov.Services.Product, prov.Presenters.Product)

	prov.Handlers.Shop = shop.NewHandler(prov.Services.Monitoring)

	return nil
}
