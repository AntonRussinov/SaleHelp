package di

import (
	"errors"
	"salehelp/internal/handler/alert"
	"salehelp/internal/handler/product"
	"salehelp/internal/handler/shop"
	shopaisettings "salehelp/internal/handler/shop_ai_settings"
	"salehelp/internal/handler/subscription"
	"salehelp/internal/handler/telegram"
)

type Handlers struct {
	Product        *product.Handler
	TelegramBot    *telegram.Handler
	Shop           *shop.Handler
	Subscription   *subscription.Handler
	Alert          *alert.Handler
	ShopAISettings *shopaisettings.Handler
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
	//prov.Handlers.Product = product.NewHandler(prov.Services.Product, prov.Presenters.Product)
	prov.Handlers.Subscription = subscription.NewHandler(prov.Services.Subscription)
	prov.Handlers.Alert = alert.NewHandler(prov.Services.AlertService)
	prov.Handlers.ShopAISettings = shopaisettings.NewHandler(prov.Repositories.ShopAISettings)
	return nil
}
