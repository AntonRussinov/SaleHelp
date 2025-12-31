package productsizesku_test

import (
	"context"
	"testing"

	"salehelp/internal/storage/pggorm/repository/productsizesku"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductSizeSKURepositoryListByProductSizeID(t *testing.T) {
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
		{"empty", false, nil, false},
		{"sql error", false, assert.AnError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			sizeID := 777
			s1 := builder.SKU()
			s2 := builder.SKU()
			s1.ProductSizeID = sizeID
			s2.ProductSizeID = sizeID

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.RowValues(s1)...)
				rows.AddRow(builder.RowValues(s2)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "product_size_skus" WHERE product_size_id = \$1`).
				WithArgs(sizeID).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.ListByProductSizeID(ctx, sizeID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if !tt.found {
				require.Len(t, got, 0)
				return
			}
			require.Len(t, got, 2)
		})
	}
}

func TestProductSizeSKURepositoryListByProductIDs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		productIDs []int64
		found      bool
		sqlErr     error
		wantErr    bool
	}{
		{"found", []int64{1, 2}, true, nil, false},
		{"empty input", []int64{}, false, nil, false},
		{"sql error", []int64{1}, false, assert.AnError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			if len(tt.productIDs) == 0 {
				got, err := repo.ListByProductIDs(ctx, tt.productIDs)
				require.NoError(t, err)
				require.Len(t, got, 0)
				return
			}

			rows := sqlmock.NewRows([]string{"sku", "product_id", "product_size_id"})
			if tt.found {
				rows.AddRow("sku1", 1, 10)
				rows.AddRow("sku2", 2, 20)
			}

			mock.ExpectQuery(`SELECT (.+) FROM product_size_skus pssk`).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.ListByProductIDs(ctx, tt.productIDs)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if !tt.found {
				require.Len(t, got, 0)
				return
			}

			require.Len(t, got, 2)
			require.Equal(t, productsizesku.SKUWithRefs{SKU: "sku1", ProductID: 1, ProductSizeID: 10}, got[0])
		})
	}
}
