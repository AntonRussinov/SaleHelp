package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type Config struct {
	TLSEnabled  bool
	TLSParamKey string
	ConnString  string
	SslRootCert string
	SslKey      string
	SslCert     string
	DriverName  string
	Pool        PoolConfig
}

var ErrUnknownDriverType = errors.New("unknown driver type")

// var DB *gorm.DB

// func InitDatabase(dsn string) error {
// 	var err error
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func SaveClientData(data []byte) error {
// 	// Сохранение данных в базу данных
// 	return nil
// }

// func ConnectWithConfig(log *slog.Logger, cfg *Config) (*sql.DB, error) {
// 	var conn *sql.DB
// 	var err error

// 	switch cfg.DriverName {
// 	case "pgx":
// 		conn, err = connectPostgresWithConfig(log, cfg)
// 	default:
// 		return nil, ErrUnknownDriverType
// 	}

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to database: %w", err)
// 	}

// 	log.Info("Database connection established")

// 	return conn, nil
// }

func ConnectGormWithConfig(log *slog.Logger, cfg *Config) (*gorm.DB, error) {
	// тут ты можешь собрать DSN с TLS параметрами, если нужно
	dsn := cfg.ConnString

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	applyPoolSettings(sqlDB, cfg.Pool)

	log.Info("Database connection established",
		slog.Int("max_open_conns", cfg.Pool.MaxOpenConns),
		slog.Int("max_idle_conns", cfg.Pool.MaxIdleConns),
		slog.Duration("conn_max_lifetime", cfg.Pool.ConnMaxLifetime),
		slog.Duration("conn_max_idle_time", cfg.Pool.ConnMaxIdleTime),
	)

	return gormDB, nil
}

func applyPoolSettings(sqlDB *sql.DB, p PoolConfig) {
	// если не задано в конфиге — ставим безопасные дефолты
	if p.MaxOpenConns <= 0 {
		p.MaxOpenConns = 30
	}
	if p.MaxIdleConns < 0 {
		p.MaxIdleConns = 0
	}
	if p.ConnMaxLifetime <= 0 {
		p.ConnMaxLifetime = time.Hour
	}
	// ConnMaxIdleTime можно оставить 0 (тогда не ограничено), но лучше 10–30m
	if p.ConnMaxIdleTime < 0 {
		p.ConnMaxIdleTime = 0
	}

	sqlDB.SetMaxOpenConns(p.MaxOpenConns)
	sqlDB.SetMaxIdleConns(p.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(p.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(p.ConnMaxIdleTime)
}
