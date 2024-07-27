package di

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"salehelp/internal/db"
	"salehelp/internal/db/gormdb"

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

	db, err := db.ConnectWithConfig(prov.log, &prov.cfg.DB)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDBInit, err)
	}

	prov.log.Info("database connection established")
	prov.Database.DB = db

	prov.log.Info("initializing GORM")

	gormDB, err := gormdb.New(prov.cfg.DB.DriverName, db, prov.cfg.GormLogLevel)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrGormInit, err)
	}

	prov.log.Info("GORM adapter initialized")
	prov.Database.GormDB = gormDB

	return nil
}
