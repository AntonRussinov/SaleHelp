package alert_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertRepositoryCreate(t *testing.T) {
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
			alert := builder.Alert()

			args := saveAlertSQLArgs(alert)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "alerts"(.+)`).
				WithArgs(args...).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			alertParam := alert

			err := repo.Create(ctx, &alertParam)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveAlertSQLArgs(a domain.Alert) []driver.Value {
	// id не передаём, он BIGSERIAL
	return []driver.Value{
		a.ShopID,
		a.ProductID,
		a.Type,
		a.Severity,
		a.Status,
		a.Message,
		a.Payload,
		a.TriggeredAt,
		a.ResolvedAt,
		a.AIRecommendationID,
	}
}
