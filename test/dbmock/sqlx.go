package dbmock

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func MockSQLX(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := MockDB(t)
	sqlxdb := sqlx.NewDb(db, "sqlmock")

	t.Cleanup(func() {
		_ = sqlxdb.Close()
	})

	return sqlxdb, mock
}

func MockExpectSQLX(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlxdb, mock := MockSQLX(t)

	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
	})

	return sqlxdb, mock
}
