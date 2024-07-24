package gormdb

import (
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultLogLevel = logger.LogLevel(4)

var ErrUnknownDriverType = errors.New("unknown driver type")

func New(driver string, conn *sql.DB) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch driver {
	case "mysql":
		dialector = mysql.New(mysql.Config{Conn: conn})
	case "pgx":
		dialector = postgres.New(postgres.Config{Conn: conn})
	default:
		return nil, ErrUnknownDriverType
	}

	config := gorm.Config{
		Logger: logger.Default.LogMode(defaultLogLevel),
	}

	db, err := gorm.Open(dialector, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm db: %w", err)
	}

	return db, nil
}
