package cron_test

import (
	"salehelp/internal/storage/cron"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCronUpdateRequest(t *testing.T) {
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

	//duration := time.Now().AddDate(0, 0, -3).Truncate(time.Minute) // 3 days ago

	// Expect the logs
	//mockLogger.On("Info", "UpdateOldPendingRequests started for model", mock.Anything).Once()
	//mockLogger.On("Info", "Old pending requests updated successfully for model.LimitRequest", mock.Anything).Once()

	// Expect the SQL queries with transaction
	// mockDB.ExpectBegin()
	// mockDB.ExpectExec(`UPDATE "limit_requests" SET "status"=\$1,"updated_at"=\$2 WHERE \(updated_at < \$3 AND status = \$4\) AND "limit_requests"."deleted_at" IS NULL`).
	// 	WithArgs("Rejected", sqlmock.AnyArg(), duration, "Pending").
	// 	WillReturnResult(sqlmock.NewResult(0, 0))
	// mockDB.ExpectCommit()

	// duration = time.Now().AddDate(0, 0, -3).Truncate(time.Minute) // 3 days ago

	// mockLogger.On("Info", "UpdateOldPendingRequests started for model", mock.Anything).Once()
	// mockLogger.On("Info", "Old pending requests updated successfully for model.QuotaRequest", mock.Anything).Once()
	// mockDB.ExpectBegin()
	// mockDB.ExpectExec(`UPDATE "quota_requests" SET "status"=\$1,"updated_at"=\$2 WHERE \(updated_at < \$3 AND status = \$4\) AND "quota_requests"."deleted_at" IS NULL`).
	// 	WithArgs("Rejected", sqlmock.AnyArg(), duration, "Pending").
	// 	WillReturnResult(sqlmock.NewResult(0, 0))
	// mockDB.ExpectCommit()

	cron.CronUpdateRequest(gormDB, mockLogger)

	// Ensure all expectations are met
	require.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestCronUpdateRequestWithError(t *testing.T) {
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
	//	duration := time.Now().AddDate(0, 0, -3).Truncate(time.Minute) // 3 days ago

	// Expect the logs
	//mockLogger.On("Info", "UpdateOldPendingRequests started for model", mock.Anything).Twice()
	//mockLogger.On("Error", "Error updating old pending requests for model: ", mock.Anything).Twice()

	// Expect the SQL queries with transaction and errors
	// mockDB.ExpectBegin()
	// mockDB.ExpectExec(`UPDATE "limit_requests" SET "status"=\$1,"updated_at"=\$2 WHERE \(updated_at < \$3 AND status = \$4\) AND "limit_requests"."deleted_at" IS NULL`).
	// 	WithArgs("Rejected", sqlmock.AnyArg(), duration, "Pending").
	// 	WillReturnError(ErrTest)
	// mockDB.ExpectRollback()

	// mockDB.ExpectBegin()
	// mockDB.ExpectExec(`UPDATE "quota_requests" SET "status"=\$1,"updated_at"=\$2 WHERE \(updated_at < \$3 AND status = \$4\) AND "quota_requests"."deleted_at" IS NULL`).
	// 	WithArgs("Rejected", sqlmock.AnyArg(), duration, "Pending").
	// 	WillReturnError(ErrTest)
	// mockDB.ExpectRollback()

	cron.CronUpdateRequest(gormDB, mockLogger)

	// Ensure all expectations are met
	require.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}
