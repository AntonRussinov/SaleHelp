package productsizesku_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductSizeSKURepositoryReplaceByProductSizeID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		deleteErr error
		insertErr error
		wantErr   bool
	}{
		{"success", nil, nil, false},
		{"delete error", assert.AnError, nil, true},
		{"insert error", nil, assert.AnError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			sizeID := 55
			skus := []string{"a", "b"}

			mock.ExpectBegin()
			mock.ExpectExec(`DELETE FROM "product_size_skus" WHERE product_size_id = \$1`).
				WithArgs(sizeID).
				WillReturnResult(sqlmock.NewResult(0, 3)).
				WillReturnError(tt.deleteErr)

			if tt.deleteErr == nil {
				mock.ExpectQuery(`INSERT INTO "product_size_skus"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).
					WillReturnError(tt.insertErr)
			}

			if tt.deleteErr == nil && tt.insertErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			err := repo.ReplaceByProductSizeID(ctx, sizeID, skus)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestProductSizeSKURepositoryUpsertMany(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		productSize int
		skus        []string
		sqlErr      error
		wantErr     bool
		expectSQL   bool
	}{
		{"success", 10, []string{"sku1", "sku2"}, nil, false, true},
		{"sql error", 10, []string{"sku1"}, assert.AnError, true, true},
		{"empty skus", 10, []string{}, nil, false, false},
		{"zero size id", 0, []string{"sku1"}, nil, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			if tt.expectSQL {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "product_size_skus"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).
					WillReturnError(tt.sqlErr)

				if tt.sqlErr == nil {
					mock.ExpectCommit()
				} else {
					mock.ExpectRollback()
				}
			}

			err := repo.UpsertMany(ctx, tt.productSize, tt.skus)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
