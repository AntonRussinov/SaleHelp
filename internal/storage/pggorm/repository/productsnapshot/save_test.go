package productsnapshot_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductSnapshotRepositoryCreate(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)
			snap := builder.Snapshot()

			args := saveProductSnapshotSQLArgs(snap)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "product_snapshots"`).
				WithArgs(args...).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			snapCopy := snap
			err := repo.Create(ctx, &snapCopy)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveProductSnapshotSQLArgs(s domain.ProductSnapshot) []driver.Value {
	return []driver.Value{
		s.ProductID,
		s.SnapshotTime,
		s.CurrencyCode,
		s.BasePrice,
		s.FinalPrice,
		s.ClubPrice,
		s.Discount,
		s.ClubDiscount,
		s.OldPrice,
		s.StockTotal,
		s.StockInWay,
		s.Rating,
		s.ReviewsCount,
		s.Sales24h,
		s.Sales7d,
		s.PositionSearch,
		s.IsBadTurnover,
		s.EditableSizePrice,
		s.Extra,
		s.CreatedAt,
	}
}
