package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"salehelp/internal/db"
	"salehelp/internal/db/gormdb"
	"salehelp/pkg/logger/sl"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	Config struct {
		Listen         string
		Environment    string
		DB             db.Config
		AllowedOrigins []string
		//GCS         gcs.Config
	}

	Service struct {
		cfg      *Config
		router   *echo.Echo
		log      *slog.Logger
		services *Services
	}
)

const shutdownTimeout = 10 * time.Second

func NewService(ctx context.Context, cfg *Config, log *slog.Logger) (*Service, error) {
	const op = "app.NewService"

	sqlDB, err := db.ConnectWithConfig(log, &cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed db connection %s: %w", op, err)
	}

	gormDB, err := gormdb.New(cfg.DB.DriverName, sqlDB)
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm connector %s: %w", op, err)
	}

	repositories := NewRepository(log, gormDB)
	services := NewServices(log, repositories)
	presenters := NewPreseneter()

	handlers, err := NewHandler(cfg, log, services, presenters)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize handlers: %w", err)
	}
	//nolint:contextcheck
	service := &Service{
		cfg:      cfg,
		log:      log,
		router:   NewRouter(cfg, log, handlers),
		services: services,
	}

	return service, nil
}

func (s *Service) Run() {
	s.log.Info("Starting application http service")

	// InitBackgroundJobs(s.log, s.services)

	if err := s.router.Start(s.cfg.Listen); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("Server closed under request")
		} else {
			s.log.Error("Server closed unexpect")
		}
	}
}

func (s *Service) Stop() error {
	if s.router == nil {
		s.log.Error("Server Close: ", sl.Err(http.ErrServerClosed))
	}

	s.log.Info("Shutdown service",
		slog.Duration("timeout", shutdownTimeout),
	)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	if err := s.router.Shutdown(ctx); err != nil {
		cancel()
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	cancel()

	return nil
}
