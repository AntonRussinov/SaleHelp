package productsizesnapshot_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductSizeSnapshotRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
	}{
		{"success", nil, false},
		{"sql error", assert.AnError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)
			s := builder.SizeSnapshot()

			args := saveArgs(s)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "product_size_snapshots"`).
				WithArgs(args...).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			param := s
			err := repo.Create(ctx, &param)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveArgs(s domain.ProductSizeSnapshot) []driver.Value {
	return []driver.Value{
		s.ProductSnapshotID,
		s.ProductSizeID,
		s.SKU,
		s.Stock,
		s.StockByWarehouse,
		s.Extra,
	}
}
