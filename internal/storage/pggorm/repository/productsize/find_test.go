package productsize_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductSizeRepositoryGetByID(t *testing.T) {
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

			ps := builder.ProductSize()
			ps.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.RowValues(ps)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "product_sizes" WHERE id = \$1`).
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
			require.Equal(t, ps.ID, got.ID)
		})
	}
}

func TestProductSizeRepositoryListByProductID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	ps1 := builder.ProductSize()
	ps2 := builder.ProductSize()
	productID := 77
	ps1.ProductID = productID
	ps2.ProductID = productID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.RowValues(ps1)...).
		AddRow(builder.RowValues(ps2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "product_sizes" WHERE product_id = \$1`).
		WithArgs(productID).
		WillReturnRows(rows)

	list, err := repo.ListByProductID(ctx, productID)

	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestProductSizeRepositoryFindByProductAndExternalID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	ps := builder.ProductSize()
	ps.ProductID = 12
	ps.ExternalID = "12345"

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.RowValues(ps)...)

	mock.ExpectQuery(`SELECT (.+) FROM "product_sizes" WHERE product_id = \$1 AND external_id = \$2`).
		WithArgs(ps.ProductID, ps.ExternalID, 1).
		WillReturnRows(rows)

	got, err := repo.FindByProductAndExternalID(ctx, ps.ProductID, ps.ExternalID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, ps.ExternalID, got.ExternalID)
}
