package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"salehelp/pkg/logger/sl"

	"github.com/jackc/pgx/v5"
)

const retriesCount = 5

func connectPostgresWithConfig(log *slog.Logger, cfg *Config) (*sql.DB, error) {
	const op = "dbconn.connectPostgresWithConfig"

	log.Debug("Initializing db connection with pgx driver")

	connStr, err := pgx.ParseConfig(cfg.ConnString)
	if err != nil {
		log.Error("failed to parse connection string", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d database=%s",
		connStr.Host, connStr.User, connStr.Password, connStr.Port, connStr.Database)

	if cfg.TLSEnabled {
		log.Debug("configuring tls connection")

		dsn += fmt.Sprintf(" sslmode=require sslrootcert=%s sslcert=%s sslkey=%s",
			cfg.SslRootCert, cfg.SslCert, cfg.SslKey)
	}

	return connectWithRetry(log, "pgx", dsn, retriesCount)
}
