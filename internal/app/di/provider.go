package di

import (
	"context"
	"log/slog"
	"net/http"
	wbclient "salehelp/internal/adapter/wb/client"
	wbdto "salehelp/internal/adapter/wb/dto"
	"salehelp/internal/config"
	"salehelp/internal/scheduler"
	"salehelp/pkg/extendedvalidator"
	"time"

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
	scheduler       *scheduler.Scheduler
	schedulerLoc    *time.Location
	telegramAdapter interface{ SendMessage(context.Context, int64, string) error }
	// AccessManager  *access.Manager
	// KaspiAPI       kaspiAPI.API
	// KaspiPortalAPI portalAPI.PortalAPI
	WBClient *wbclient.Client
}

// TelegramAdapter возвращает адаптер для отправки сообщений в Telegram (для cron-задач).
func (prov *Provider) TelegramAdapter() interface{ SendMessage(context.Context, int64, string) error } {
	return prov.telegramAdapter
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

	// Часовой пояс для отчётов WB (sales, orders): в каком поясе WB отдаёт даты и ждёт dateFrom.
	reportTZ := prov.cfg.WBReportTimezone
	if reportTZ == "" {
		reportTZ = "Europe/Moscow"
	}
	loc, err := time.LoadLocation(reportTZ)
	if err != nil {
		prov.log.Warn("wb report timezone load failed, using Moscow", "tz", reportTZ, "err", err)
		loc, _ = time.LoadLocation("Europe/Moscow")
		if loc == nil {
			loc = time.FixedZone("MSK", 3*3600)
		}
	}
	wbdto.SetWBReportLocation(loc)
	prov.WBClient = wbclient.New(httpClient, loc)
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
