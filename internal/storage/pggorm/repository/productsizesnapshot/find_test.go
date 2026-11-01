package productsizesnapshot_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductSizeSnapshotRepositoryListBySnapshotID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	s1 := builder.SizeSnapshot()
	s2 := builder.SizeSnapshot()

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.RowValues(s1)...).
		AddRow(builder.RowValues(s2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "product_size_snapshots" WHERE product_snapshot_id = \$1`).
		WithArgs(s1.ProductSnapshotID).
		WillReturnRows(rows)

	list, err := repo.ListBySnapshotID(ctx, s1.ProductSnapshotID)

	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestProductSizeSnapshotRepositoryGetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		found   bool
		sqlErr  error
		wantErr bool
	}{
		{"found", true, nil, false},
		{"not found", false, nil, false},
		{"sql error", false, assert.AnError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			s := builder.SizeSnapshot()
			s.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.RowValues(s)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "product_size_snapshots" WHERE id = \$1`).
				WithArgs(1, 1).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, 1)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if !tt.found {
				require.Nil(t, got)
				return
			}

			require.Equal(t, s.ID, got.ID)
		})
	}
}
