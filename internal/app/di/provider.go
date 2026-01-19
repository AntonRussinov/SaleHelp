package di

import (
	"context"
	"log/slog"
	"net/http"
	wbclient "salehelp/internal/adapter/wb/client"
	"salehelp/internal/config"
	"salehelp/internal/scheduler"
	"salehelp/pkg/extendedvalidator"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
)

type Provider struct {
	cfg          *config.ServiceConfig
	log          *slog.Logger
	Cron         *cron.Cron
	Database     Database
	Repositories Repositories
	Services     Services
	Presenters   Presenters
	Handlers     Handlers
	Router       *echo.Echo
	Validator    *extendedvalidator.Validator
	scheduler    *scheduler.Scheduler
	// AccessManager  *access.Manager
	// KaspiAPI       kaspiAPI.API
	// KaspiPortalAPI portalAPI.PortalAPI
	WBClient *wbclient.Client
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

func (prov *Provider) initWBClient() {
	httpClient := &http.Client{}

	// Создаём WB-клиент (New возвращает *Client)
	prov.WBClient = wbclient.New(httpClient)
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

	prov.initWBClient()

	prov.InitRepositories()
	prov.InitServices()
	prov.InitPresenters()

	// prov.InitKaspiPortalAPI(prov.Repositories.Merchant)
	// err = prov.InitCronRunner(ctx)
	// if err != nil {
	// 	return err
	// }

	err = prov.InitHandlers()
	if err != nil {
		return err
	}

	err = prov.InitValidator()
	if err != nil {
		return err
	}

	err = prov.initScheduler(ctx)
	if err != nil {
		return err
	}

	prov.InitRouter()

	return nil
}
