package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	TLSEnabled  bool
	TLSParamKey string
	ConnString  string
	SslRootCert string
	SslKey      string
	SslCert     string
	DriverName  string
}

var ErrUnknownDriverType = errors.New("unknown driver type")

var DB *gorm.DB

func InitDatabase(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func SaveClientData(data []byte) error {
	// Сохранение данных в базу данных
	return nil
}

func ConnectWithConfig(log *slog.Logger, cfg *Config) (*sql.DB, error) {
	var conn *sql.DB
	var err error

	switch cfg.DriverName {
	case "pgx":
		conn, err = connectPostgresWithConfig(log, cfg)
	default:
		return nil, ErrUnknownDriverType
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info("Database connection established")

	return conn, nil
}
