package dbmock

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MockGORM(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := MockDB(t)
	dialector := postgres.New(postgres.Config{Conn: db})
	gormconfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	gormdb, err := gorm.Open(dialector, gormconfig)
	require.NoError(t, err)

	return gormdb, mock
}

func MockExpectGORM(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	gormdb, mock := MockGORM(t)

	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
	})

	return gormdb, mock
}
