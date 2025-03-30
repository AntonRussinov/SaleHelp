package stock_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/test/dbmock"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStockRepositorySave(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		stock   domain.Stock
		sqlErr  error
		wantErr bool
	}{
		{name: "executes insert query", stock: builder.Stock()},
		{name: "returns error on sql insert error", stock: builder.Stock(), sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "stocks"(.+)`).WithArgs(saveStockSQLArgs(tt.stock)...).
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			stockParam := tt.stock

			err := repo.Save(ctx, &stockParam)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveStockSQLArgs(stock domain.Stock) []driver.Value {
	values := []driver.Value{

		stock.MerchantUID,
		stock.OfferID,
		stock.CityID,
		stock.StoreID,
		stock.Available,
		stock.PreOrder,
		stock.StockCount,
		stock.Enabled,
		dbmock.MatchTime(),
		dbmock.MatchTime(),
		nil,
	}

	return values
}
