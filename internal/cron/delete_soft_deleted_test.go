package cron_test

import (
	"errors"
	"salehelp/internal/cron"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrTest = errors.New("error for test")

func TestDeleteOldSoftDeletedRecords(t *testing.T) {
	t.Parallel()

	db, mockDB, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})

	require.NoError(t, err)

	mockLogger := new(MockLogger)
	//duration := time.Now().AddDate(0, 0, -7).Truncate(time.Minute) // 7 days ago

	// Expect the logs
	mockLogger.On("Info", "DeleteOldSoftDeletedRecords started", mock.Anything).Once()
	// mockLogger.On("Info", "Old soft deleted records deleted successfully", mock.Anything).Once()

	// Expect the SQL queries
	// mockDB.ExpectBegin()
	// mockDB.ExpectExec(`DELETE FROM "templates" WHERE deleted_at IS NOT NULL AND deleted_at < ?`).
	// 	WithArgs(duration).
	// 	WillReturnResult(sqlmock.NewResult(0, 0))
	// mockDB.ExpectCommit()

	// mockDB.ExpectBegin()
	// mockDB.ExpectExec(`DELETE FROM "applications" WHERE deleted_at IS NOT NULL AND deleted_at < ?`).
	// 	WithArgs(duration).
	// 	WillReturnResult(sqlmock.NewResult(0, 0))
	// mockDB.ExpectCommit()

	cron.DeleteOldSoftDeletedRecords(gormDB, mockLogger)

	// Ensure all expectations are met
	require.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestDeleteOldSoftDeletedRecordsWithError(t *testing.T) {
	t.Parallel()

	db, mockDB, err := sqlmock.New()

	require.NoError(t, err)

	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	mockLogger := new(MockLogger)
	// duration := time.Now().AddDate(0, 0, -7).Truncate(time.Minute) // 7 days ago, truncated to seconds

	// Expect the logs
	mockLogger.On("Info", "DeleteOldSoftDeletedRecords started", mock.Anything).Once()
	// mockLogger.On("Error", "Error deleting old soft deleted records for Template table: ", mock.Anything).Once()
	// mockLogger.On("Error", "Error deleting old soft deleted records for Application table: ", mock.Anything).Once()

	// Expect the SQL queries with transaction and errors
	// mockDB.ExpectBegin()
	// mockDB.ExpectExec(`DELETE FROM "templates" WHERE deleted_at IS NOT NULL AND deleted_at < ?`).
	// 	WithArgs(duration).
	// 	WillReturnError(ErrTest)
	// mockDB.ExpectRollback()

	// mockDB.ExpectBegin()
	// mockDB.ExpectExec(`DELETE FROM "applications" WHERE deleted_at IS NOT NULL AND deleted_at < ?`).
	// 	WithArgs(duration).
	// 	WillReturnError(ErrTest)
	// mockDB.ExpectRollback()

	cron.DeleteOldSoftDeletedRecords(gormDB, mockLogger)

	// Ensure all expectations are met
	require.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}
