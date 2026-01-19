package alerttriggered_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertTriggeredRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
	}{
		{name: "executes insert query"},
		{name: "returns error on sql insert error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)
			a := builder.AlertTriggered()

			args := saveAlertTriggeredSQLArgs(a)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "alerts_triggered"(.+)`).
				WithArgs(args...).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			param := a

			err := repo.Create(ctx, &param)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveAlertTriggeredSQLArgs(a domain.AlertTriggered) []driver.Value {
	// id — BIGSERIAL, не передаём
	return []driver.Value{
		a.AlertID,
		a.ProductID,
		a.SnapshotID,
		a.ShopID,
		a.Message,
		a.ProductSizeID,
		a.AlertType,
		a.CreatedAt,
	}
}
