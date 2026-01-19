package sellerwarehouse_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSellerWarehouseRepositoryGetByShopAndWBOfficeID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		found   bool
		sqlErr  error
		wantErr bool
	}{
		{name: "found", found: true},
		{name: "not found", found: false},
		{name: "sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			w := builder.Warehouse()
			w.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.WarehouseRowValues(w)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "seller_warehouses" WHERE shop_id = \$1 AND wb_office_id = \$2`).
				WithArgs(w.ShopID, w.WBOfficeID).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByShopAndWBOfficeID(ctx, w.ShopID, w.WBOfficeID)

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
			require.Equal(t, w.ShopID, got.ShopID)
			require.Equal(t, w.WBOfficeID, got.WBOfficeID)
		})
	}
}

func TestSellerWarehouseRepositoryListByShop(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	w1 := builder.Warehouse()
	w2 := builder.Warehouse()
	shopID := w1.ShopID
	w2.ShopID = shopID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.WarehouseRowValues(w1)...).
		AddRow(builder.WarehouseRowValues(w2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "seller_warehouses" WHERE shop_id = \$1`).
		WithArgs(shopID).
		WillReturnRows(rows)

	list, err := repo.ListByShop(ctx, shopID)

	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, shopID, list[0].ShopID)
	require.Equal(t, shopID, list[1].ShopID)
}
