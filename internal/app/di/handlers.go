package di

import (
	"errors"
	"salehelp/internal/handler/product"
	"salehelp/internal/handler/shop"
	"salehelp/internal/handler/telegram"
)

type Handlers struct {
	Product     *product.Handler
	TelegramBot *telegram.Handler
	Shop        *shop.Handler
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
	prov.Handlers.TelegramBot = telegram.NewHandler(prov.Services.TelegramBot, prov.cfg.TelegramSecret)

	return nil
}
