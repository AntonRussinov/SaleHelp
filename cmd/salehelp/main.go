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

	log.Println("reading config file path from command line arguments")

	configFile := pflag.StringP("config", "c", "configs/config.yml", "path to configuration file")

	pflag.Parse()

	log.Println("using config file ", *configFile)
	// Загрузка конфигурации
	cfg, err := config.Read(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("initializing service")

	applogger := logger.New(cfg.Environment)
	applogger.Info("initializing service",
		slog.String("config", *configFile),
		slog.String("env", cfg.Environment),
		slog.String("addr", cfg.Service.Listen),
	)

	go shutdownHandler(ctx, cancel)

	// service, err := app.NewService(ctx, &cfg.Service, applogger)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// go service.Run()

	// <-ctx.Done()

	// err = service.Stop()
	// if err != nil {
	// 	applogger.Error("Server Close: ", sl.Err(err))
	// }

	// // Инициализация базы данных
	// // err = db.InitDatabase(cfg.Service.DB)
	// // if err != nil {
	// // 	log.Fatalf("Failed to connect to database: %v", err)
	// // }

	// // Запуск планировщика задач
	// scheduler.StartScheduler()

	// // Блокировка основного потока для предотвращения завершения приложения
	// select {}

	prov := di.NewProvider(applogger, cfg.Service)
	application := app.New(cfg.Service, applogger, prov)

	err = application.Run(ctx)
	if err != nil {
		applogger.Error("failed to run application", slog.Any("error", err))
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	err = application.Shutdown(shutdownCtx)
	if err != nil {
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
