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

		StartScheduler(ctx context.Context) error
		StopScheduler(ctx context.Context) error
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

	if err := app.prov.StartScheduler(ctx); err != nil {
		return err
	}

	// старт сервера в горутине
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.prov.GetRouter().Start(app.cfg.Listen)
	}()

	// ждём либо ctx.Done(), либо ошибку сервера
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		// сервер упал сам
		if errors.Is(err, http.ErrServerClosed) {
			app.log.Info("Server closed")
			return nil
		}
		app.log.Error("Server stopped unexpectedly", slog.Any("error", err))
		return err
	}
}

func (app *APP) Shutdown(ctx context.Context) error {
	if app.isShuttingDown() {
		return ErrAppStopped
	}

	app.inShutdown.Store(true)

	// graceful shutdown HTTP сервера
	if err := app.prov.GetRouter().Shutdown(ctx); err != nil {
		return err
	}

	if err := app.prov.StopScheduler(ctx); err != nil {
		return err
	}

	// если есть ресурсы: db, clients — можно закрывать тут
	// _ = app.prov.Close()

	return nil
}

func (app *APP) isShuttingDown() bool {
	return app.inShutdown.Load()
}
