package wbwarehouse_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWBWarehouseRepositoryGetByWBID(t *testing.T) {
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

			w := builder.Warehouse()
			w.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.WarehouseRowValues(w)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "wb_warehouses" WHERE wb_id = \$1`).
				WithArgs(w.WBID, 1).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByWBID(ctx, w.WBID)

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
			require.Equal(t, w.WBID, got.WBID)
			require.Equal(t, w.Name, got.Name)
		})
	}
}

func TestWBWarehouseRepositoryListAll(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	w1 := builder.Warehouse()
	w2 := builder.Warehouse()

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.WarehouseRowValues(w1)...).
		AddRow(builder.WarehouseRowValues(w2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "wb_warehouses" ORDER BY id`).
		WillReturnRows(rows)

	list, err := repo.ListAll(ctx)

	require.NoError(t, err)
	require.Len(t, list, 2)
}
