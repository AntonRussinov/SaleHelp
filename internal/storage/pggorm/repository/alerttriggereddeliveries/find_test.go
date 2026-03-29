package alerttriggereddeliveries_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertTriggeredDeliveryRepositoryGetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		found   bool
		sqlErr  error
		wantErr bool
	}{
		{name: "found", found: true},
		{name: "not found", found: false},
		{name: "sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			d := builder.Delivery()
			d.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.DeliveryRowValues(d)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "alerts_triggered_deliveries" WHERE id = \$1`).
				WithArgs(d.ID, 1).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, d.ID)

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
			require.Equal(t, d.ID, got.ID)
			require.Equal(t, d.TriggeredID, got.TriggeredID)
		})
	}
}

func TestAlertTriggeredDeliveryRepositoryListByTriggeredID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	d1 := builder.Delivery()
	d2 := builder.Delivery()
	triggeredID := d1.TriggeredID
	d2.TriggeredID = triggeredID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.DeliveryRowValues(d1)...).
		AddRow(builder.DeliveryRowValues(d2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "alerts_triggered_deliveries" WHERE triggered_id = \$1 ORDER BY id`).
		WithArgs(triggeredID).
		WillReturnRows(rows)

	list, err := repo.ListByTriggeredID(ctx, triggeredID)

	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestAlertTriggeredDeliveryRepositoryListReady(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	d1 := builder.Delivery()
	d2 := builder.Delivery()

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.DeliveryRowValues(d1)...).
		AddRow(builder.DeliveryRowValues(d2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "alerts_triggered_deliveries" WHERE status IN \('new','failed'\) AND \(lease_until IS NULL OR lease_until < now\(\)\) ORDER BY id LIMIT \$1`).
		WithArgs(50).
		WillReturnRows(rows)

	list, err := repo.ListReady(ctx, 50)

	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestAlertTriggeredDeliveryRepositoryLeaseBatch(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
	}{
		{name: "success"},
		{name: "sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			leaseOwner := "worker-1"
			leaseFor := 30 * time.Second
			limit := 2

			d1 := builder.Delivery()
			d2 := builder.Delivery()

			rows := sqlmock.NewRows(JoinedColumns()).
				AddRow(builder.DeliveryRowValues(d1)...).
				AddRow(builder.DeliveryRowValues(d2)...)

			mock.ExpectQuery(`WITH cte AS \(`).
				// тут WithArgs можно, но leaseUntil = time.Now().Add(...) не стабилен,
				// поэтому args не проверяем
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.LeaseBatch(ctx, leaseOwner, leaseFor, limit)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, got, 2)
		})
	}
}
