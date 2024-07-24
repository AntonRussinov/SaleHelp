package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"salehelp/pkg/logger/sl"
	"time"
)

const retryTimeout = 5 * time.Second

func connect(driver, dsn string) (*sql.DB, error) {
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return conn, nil
}

func connectWithRetry(log *slog.Logger, driver string, dsn string, n int) (conn *sql.DB, err error) {
	log.Info("trying to connect to db", slog.Int("retries", n))

	log.Debug("connecting to db", slog.String("dsn", dsn))

	for range n {
		conn, err = connect(driver, dsn)
		if err != nil {
			log.Error("failed to connect to db", sl.Err(err))
			time.Sleep(retryTimeout)

			continue
		}

		return conn, err
	}

	return nil, fmt.Errorf("failed to connect to database: %w", err)
}
