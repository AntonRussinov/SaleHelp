package product_test

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

func TestProductRepositorySave(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		product domain.Product
		sqlErr  error
		wantErr bool
	}{
		{name: "executes insert query", product: builder.Product()},
		{name: "returns error on sql insert error", product: builder.Product(), sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "products"(.+)`).WithArgs(saveProductSQLArgs(tt.product)...).
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			productParam := tt.product

			err := repo.Save(ctx, &productParam)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveProductSQLArgs(product domain.Product) []driver.Value {
	values := []driver.Value{
		product.ID,

		product.MasterSKU,
		product.Model,
		product.Brand,
		product.MasterTitle,
		product.Title,
		product.VerticalCategory,
		product.MasterCategory,
		dbmock.MatchTime(),
		dbmock.MatchTime(),
		nil,
	}

	return values
}
