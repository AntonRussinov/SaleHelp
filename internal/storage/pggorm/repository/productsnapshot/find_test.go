package productsnapshot_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductSnapshotRepositoryGetLastByProductID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name      string
		found     bool
		sqlErr    error
		wantErr   bool
		productID int
	}{
		{name: "found", found: true, productID: 10},
		{name: "not found", found: false, productID: 10},
		{name: "sql error", sqlErr: assert.AnError, wantErr: true, productID: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			snap := builder.Snapshot()
			snap.ProductID = tt.productID

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.ProductSnapshotRowValues(snap)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "product_snapshots" WHERE product_id = \$1 ORDER BY snapshot_time DESC LIMIT 1`).
				WithArgs(tt.productID).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetLastByProductID(ctx, tt.productID)

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
			require.Equal(t, snap.ProductID, got.ProductID)
		})
	}
}

func TestProductSnapshotRepositoryListByProductID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	productID := 10

	s1 := builder.Snapshot()
	s1.ProductID = productID

	s2 := builder.Snapshot()
	s2.ProductID = productID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.ProductSnapshotRowValues(s1)...).
		AddRow(builder.ProductSnapshotRowValues(s2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "product_snapshots" WHERE product_id = \$1 ORDER BY snapshot_time DESC`).
		WithArgs(productID).
		WillReturnRows(rows)

	list, err := repo.ListByProductID(ctx, productID, 0, 0)

	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, productID, list[0].ProductID)
	require.Equal(t, productID, list[1].ProductID)
}
