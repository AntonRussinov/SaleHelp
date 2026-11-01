// internal/storage/pggorm/repository/shop/get_by_id_test.go
package shop_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShopRepositoryGetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
		found   bool
	}{
		{name: "returns shop", found: true},
		{name: "returns nil when not found", found: false},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			shop := builder.Shop()
			shop.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.ShopRowValues(shop)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "shops" WHERE id = \$1(.+)`).
				WithArgs(shop.ID, 1).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, shop.ID)

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
			require.Equal(t, shop.ID, got.ID)
			require.Equal(t, shop.UserID, got.UserID)
			require.Equal(t, shop.Marketplace, got.Marketplace)
		})
	}
}

func TestShopRepositoryGetByUserAndMarketplace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name        string
		sqlErr      error
		wantErr     bool
		found       bool
		userID      int
		marketplace string
	}{
		{name: "returns shop", found: true, userID: 10, marketplace: "WB"},
		{name: "returns nil when not found", found: false, userID: 10, marketplace: "WB"},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true, userID: 10, marketplace: "WB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			shop := builder.Shop()
			shop.UserID = tt.userID
			shop.Marketplace = tt.marketplace

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.ShopRowValues(shop)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "shops" WHERE user_id = \$1 AND marketplace = \$2(.+)`).
				WithArgs(tt.userID, tt.marketplace, 1).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByUserAndMarketplace(ctx, tt.userID, tt.marketplace)

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
			require.Equal(t, shop.UserID, got.UserID)
			require.Equal(t, shop.Marketplace, got.Marketplace)
		})
	}
}

func TestShopRepositoryListByUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	userID := 42
	s1 := builder.Shop()
	s1.UserID = userID
	s2 := builder.Shop()
	s2.UserID = userID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.ShopRowValues(s1)...).
		AddRow(builder.ShopRowValues(s2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "shops" WHERE user_id = \$1(.+)`).
		WithArgs(userID).
		WillReturnRows(rows)

	shops, err := repo.ListByUser(ctx, userID, 0, 0)
	require.NoError(t, err)
	require.Len(t, shops, 2)

	require.Equal(t, s1.UserID, shops[0].UserID)
	require.Equal(t, s2.UserID, shops[1].UserID)
}
