package product_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductRepositoryGetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
		found   bool
	}{
		{name: "returns product", found: true},
		{name: "returns nil when not found", found: false},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			product := builder.Product()
			product.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.ProductRowValues(product)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "products" WHERE id = \$1(.+)`).
				WithArgs(product.ID, 1).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, product.ID)

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
			require.Equal(t, product.ID, got.ID)
			require.Equal(t, product.ShopID, got.ShopID)
			require.Equal(t, product.ExternalID, got.ExternalID)
		})
	}
}

func TestProductRepositoryGetByShopAndExternalID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name       string
		sqlErr     error
		wantErr    bool
		found      bool
		shopID     int
		externalID string
	}{
		{name: "returns product", found: true, shopID: 10, externalID: "ext-123"},
		{name: "returns nil when not found", found: false, shopID: 10, externalID: "ext-123"},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true, shopID: 10, externalID: "ext-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			product := builder.Product()
			product.ShopID = tt.shopID
			product.ExternalID = tt.externalID

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.ProductRowValues(product)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "products" WHERE shop_id = \$1 AND external_id = \$2(.+)`).
				WithArgs(tt.shopID, tt.externalID, 1).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByShopAndExternalID(ctx, tt.shopID, tt.externalID)

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
			require.Equal(t, product.ShopID, got.ShopID)
			require.Equal(t, product.ExternalID, got.ExternalID)
		})
	}
}

func TestProductRepositoryListByShop(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	shopID := 42
	p1 := builder.Product()
	p1.ShopID = shopID
	p2 := builder.Product()
	p2.ShopID = shopID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.ProductRowValues(p1)...).
		AddRow(builder.ProductRowValues(p2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "products" WHERE shop_id = \$1(.+)`).
		WithArgs(shopID).
		WillReturnRows(rows)

	products, err := repo.ListByShop(ctx, shopID, 0, 0)
	require.NoError(t, err)
	require.Len(t, products, 2)

	require.Equal(t, p1.ShopID, products[0].ShopID)
	require.Equal(t, p2.ShopID, products[1].ShopID)
}
