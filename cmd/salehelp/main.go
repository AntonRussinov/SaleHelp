package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"salehelp/internal/app"
	"salehelp/internal/config"
	"salehelp/internal/scheduler"
	"salehelp/pkg/logger"
	"salehelp/pkg/logger/sl"
	"syscall"

	"github.com/spf13/pflag"
)

func main() {
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	service, err := app.NewService(ctx, &cfg.Service, applogger)
	if err != nil {
		log.Panic(err)
	}

	go service.Run()

	<-ctx.Done()

	err = service.Stop()
	if err != nil {
		applogger.Error("Server Close: ", sl.Err(err))
	}

	// Инициализация базы данных
	// err = db.InitDatabase(cfg.Service.DB)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }

	// Запуск планировщика задач
	scheduler.StartScheduler()

	// Блокировка основного потока для предотвращения завершения приложения
	select {}
}
