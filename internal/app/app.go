package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"salehelp/internal/config"
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

type (
	Provider interface {
		Init(ctx context.Context) error
		GetRouter() *echo.Echo
	}

	APP struct {
		prov       Provider
		cfg        *config.ServiceConfig
		log        *slog.Logger
		inShutdown atomic.Bool
	}
)

var ErrAppStopped = errors.New("app is stopped")

func New(cfg *config.ServiceConfig, log *slog.Logger, prov Provider) *APP {
	return &APP{
		cfg:  cfg,
		log:  log,
		prov: prov,
	}
}

func (app *APP) Run(ctx context.Context) error {
	if app.isShuttingDown() {
		return ErrAppStopped
	}

	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := app.prov.Init(appCtx)
	if err != nil {
		return err
	}

	if err := app.prov.GetRouter().Start(app.cfg.Listen); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			app.log.Info("Server closed under request")
		} else {
			app.log.Error("Server closed unexpect")
		}
	}

	<-appCtx.Done()

	return appCtx.Err()
}

func (app *APP) Shutdown(ctx context.Context) error {
	if app.isShuttingDown() {
		return ErrAppStopped
	}

	app.inShutdown.Store(true)

	return nil
}

func (app *APP) isShuttingDown() bool {
	return app.inShutdown.Load()
}
