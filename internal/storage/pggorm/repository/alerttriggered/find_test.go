package alerttriggered_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertTriggeredRepositoryGetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		found   bool
		sqlErr  error
		wantErr bool
	}{
		{name: "returns record", found: true},
		{name: "returns nil when not found", found: false},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			a := builder.AlertTriggered()
			a.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.AlertTriggeredRowValues(a)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "alerts_triggered" WHERE id = \$1(.+)`).
				WithArgs(a.ID).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, a.ID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if !tt.found {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, a.ID, got.ID)
			require.Equal(t, a.AlertID, got.AlertID)
		})
	}
}

func TestAlertTriggeredRepositoryListByAlertID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	alertID := 10

	a1 := builder.AlertTriggered()
	a1.AlertID = alertID

	a2 := builder.AlertTriggered()
	a2.AlertID = alertID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.AlertTriggeredRowValues(a1)...).
		AddRow(builder.AlertTriggeredRowValues(a2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "alerts_triggered" WHERE alert_id = \$1(.+)`).
		WithArgs(alertID).
		WillReturnRows(rows)

	list, err := repo.ListByAlertID(ctx, alertID, 0, 0)

	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, alertID, list[0].AlertID)
	require.Equal(t, alertID, list[1].AlertID)
}

func TestAlertTriggeredRepositoryListByProductID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	productID := 777

	a1 := builder.AlertTriggered()
	a1.ProductID = productID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.AlertTriggeredRowValues(a1)...)

	mock.ExpectQuery(`SELECT (.+) FROM "alerts_triggered" WHERE product_id = \$1(.+)`).
		WithArgs(productID).
		WillReturnRows(rows)

	list, err := repo.ListByProductID(ctx, productID, 0, 0)

	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, productID, list[0].ProductID)
}
