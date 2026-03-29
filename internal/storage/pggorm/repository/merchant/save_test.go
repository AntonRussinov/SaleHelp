package merchant_test

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

func TestMerchantRepositorySave(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	builder := NewDataBuilder(t)

	tests := []struct {
		name     string
		merchant domain.Merchant
		sqlErr   error
		wantErr  bool
	}{
		{name: "executes insert query", merchant: builder.Merchant()},
		{name: "returns error on sql insert error", merchant: builder.Merchant(), sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "merchants"(.+)`).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			merchantParam := tt.merchant

			err := repo.Save(ctx, &merchantParam)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveMerchantSQLArgs(merchant domain.Merchant) []driver.Value {
	values := []driver.Value{
		merchant.MerchantUID,
		merchant.CompanyName,
		merchant.IsEnabled,
		dbmock.MatchTime(),
		dbmock.MatchTime(),
		nil,
	}

	return values
}
