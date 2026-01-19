package di

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"salehelp/internal/db"

	"gorm.io/gorm"
)

type Database struct {
	DB     *sql.DB
	GormDB *gorm.DB
}

var (
	ErrDBInit   = errors.New("failed to initialize database connection")
	ErrGormInit = errors.New("failed to initialize GORM")
)

func (prov *Provider) InitDB(ctx context.Context) error {
	prov.log.Info("initializing database connection...")

	// 1. Создаём GORM + настраиваем пул внутри
	gormDB, err := db.ConnectGormWithConfig(prov.log, &prov.cfg.DB)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDBInit, err)
	}

	// 2. Получаем underlying *sql.DB
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}

	prov.Database.GormDB = gormDB
	prov.Database.DB = sqlDB

	prov.log.Info("database initialized successfully")

	return nil
}
