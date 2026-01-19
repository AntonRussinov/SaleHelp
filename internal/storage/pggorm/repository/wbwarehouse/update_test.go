package wbwarehouse_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWBWarehouseRepositoryUpdate(t *testing.T) {
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
			w := builder.Warehouse()
			w.ID = 1

			mock.ExpectBegin()
			mock.ExpectExec(`UPDATE "wb_warehouses" SET`).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			param := w
			err := repo.Update(ctx, &param)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestWBWarehouseRepositoryUpsert(t *testing.T) {
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
			w := builder.Warehouse()

			mock.ExpectBegin()
			// GORM с OnConflict всё равно делает INSERT ... ON CONFLICT,
			// нам достаточно матча по началу
			mock.ExpectExec(`INSERT INTO "wb_warehouses"`).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			param := w
			err := repo.Upsert(ctx, &param)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
