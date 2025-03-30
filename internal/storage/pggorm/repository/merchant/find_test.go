package merchant_test

import (
	"context"
	"salehelp/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerchantRepositoryFindByID(t *testing.T) {
	t.Parallel()

	builder := NewDataBuilder(t)

	tests := []struct {
		name     string
		merchant domain.Merchant
		sqlErr   error
		wantErr  bool
	}{
		{name: "returns merchant with matching id", merchant: builder.Merchant()},
		{name: "returns error on sql select error", merchant: builder.Merchant(), sqlErr: assert.AnError, wantErr: true},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			mock.ExpectQuery(`SELECT(.+)FROM "merchants"(.+)`).WithArgs(tt.merchant.MerchantUID, 1).
				WillReturnRows(sqlmock.NewRows(JoinedColumns()).
					AddRow(builder.MerchantRowValues(tt.merchant)...)).
				WillReturnError(tt.sqlErr)

			req, err := repo.FindByID(ctx, tt.merchant.MerchantUID)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, req)
			} else {
				require.NoError(t, err)
				require.NotNil(t, req)
				require.Equal(t, tt.merchant, *req)
			}
		})
	}
}
