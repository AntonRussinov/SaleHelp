package di

import (
	"context"
	"log/slog"
	"salehelp/internal/access"
	kaspiAPI "salehelp/internal/adapter/kaspi/api/connector"
	portalAPI "salehelp/internal/adapter/kaspi/portal/connector"
	"salehelp/internal/config"
	"salehelp/pkg/extendedvalidator"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
)

type Provider struct {
	cfg            *config.ServiceConfig
	log            *slog.Logger
	Cron           *cron.Cron
	Database       Database
	Repositories   Repositories
	Services       Services
	Presenters     Presenters
	Handlers       Handlers
	Router         *echo.Echo
	Validator      *extendedvalidator.Validator
	AccessManager  *access.Manager
	KaspiAPI       kaspiAPI.API
	KaspiPortalAPI portalAPI.PortalAPI
}

func NewProvider(log *slog.Logger, cfg *config.ServiceConfig) *Provider {
	return &Provider{
		cfg: cfg,
		log: log,
	}
}

func (prov *Provider) GetRouter() *echo.Echo {
	return prov.Router
}

func (prov *Provider) GetCron() *cron.Cron {
	return prov.Cron
}

func (prov *Provider) Init(ctx context.Context) error {
	err := prov.InitDB(ctx)
	if err != nil {
		return err
	}

	// err = prov.InitAccessManager(ctx)
	// if err != nil {
	// 	return err
	// }

	prov.InitRepositories()
	prov.InitServices()
	prov.InitPresenters()
	prov.InitKaspiAPI()
	prov.InitKaspiPortalAPI(prov.Repositories.Merchant)
	err = prov.InitCronRunner(ctx)
	if err != nil {
		return err
	}

	err = prov.InitHandlers()
	if err != nil {
		return err
	}

	err = prov.InitValidator()
	if err != nil {
		return err
	}

	prov.InitRouter()

	return nil
}
