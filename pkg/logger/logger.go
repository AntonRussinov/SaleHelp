package logger

import (
	"log/slog"
	"os"
	"salehelp/pkg/logger/handlers/prettylog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func New(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(prettylog.NewHandler(
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		))
	default:
		log = defaultLogger()
	}

	return log
}

func defaultLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
}
