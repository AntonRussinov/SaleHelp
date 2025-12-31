package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"salehelp/internal/app"
	"salehelp/internal/app/di"
	"salehelp/internal/config"
	"salehelp/pkg/logger"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

const shutdownTimeout = 15 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configFile := pflag.StringP("config", "c", "configs/config.yml", "path to configuration file")
	pflag.Parse()

	cfg, err := config.Read(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	applogger := logger.New(cfg.Environment)

	go shutdownHandler(ctx, cancel) // по сигналу вызовет cancel()

	prov := di.NewProvider(applogger, cfg.Service)
	application := app.New(cfg.Service, applogger, prov)

	// init DI / router / db
	if err := prov.Init(ctx); err != nil {
		applogger.Error("provider init failed", slog.Any("error", err))
		return
	}

	// 2) запускаем сервер в отдельной горутине (иначе дальше не дойдём)
	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- application.Run(ctx)
	}()

	// 1) запускаем воркеры СРАЗУ (они крутятся пока ctx жив)
	go prov.Services.Notifier.Run(ctx)    // deliveries -> telegram
	go prov.Services.AlertWorker.Run(ctx) // shops -> engine -> triggered+fanout

	// 3) ждём либо:
	// - сигнал (ctx.Done())
	// - либо ошибка старта/работы сервера
	select {
	case <-ctx.Done():
		applogger.Info("shutdown requested")
	case err := <-runErrCh:
		if err != nil {
			applogger.Error("application stopped with error", slog.Any("error", err))
		}
		// чтобы точно завершиться — отменим контекст
		cancel()
	}

	// 4) graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		applogger.Error("failed to shutdown application", slog.Any("error", err))
	}
}

func shutdownHandler(ctx context.Context, cancel context.CancelFunc) {
	shutdownSignals := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM}

	notify, stop := signal.NotifyContext(ctx, shutdownSignals...)
	defer stop()

	<-notify.Done()
	cancel()
}
