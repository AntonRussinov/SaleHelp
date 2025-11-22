package alert_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertRepositoryGetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		found   bool
		sqlErr  error
		wantErr bool
	}{
		{name: "returns alert", found: true},
		{name: "returns nil when not found", found: false},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			alert := builder.Alert()
			alert.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.AlertRowValues(alert)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "alerts" WHERE id = \$1(.+)`).
				WithArgs(alert.ID).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, alert.ID)

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
			require.Equal(t, alert.ID, got.ID)
			require.Equal(t, alert.ShopID, got.ShopID)
			require.Equal(t, alert.Type, got.Type)
		})
	}
}

func TestAlertRepositoryListByShop(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	shopID := 42

	a1 := builder.Alert()
	a1.ShopID = shopID
	a2 := builder.Alert()
	a2.ShopID = shopID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.AlertRowValues(a1)...).
		AddRow(builder.AlertRowValues(a2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "alerts" WHERE shop_id = \$1(.+)`).
		WithArgs(shopID).
		WillReturnRows(rows)

	list, err := repo.ListByShop(ctx, shopID, 0, 0)

	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, shopID, list[0].ShopID)
	require.Equal(t, shopID, list[1].ShopID)
}

func TestAlertRepositoryListByProduct(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	productID := 777

	a1 := builder.Alert()
	a1.ProductID = &productID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.AlertRowValues(a1)...)

	mock.ExpectQuery(`SELECT (.+) FROM "alerts" WHERE product_id = \$1(.+)`).
		WithArgs(productID).
		WillReturnRows(rows)

	list, err := repo.ListByProduct(ctx, productID, 0, 0)

	require.NoError(t, err)
	require.Len(t, list, 1)
	require.NotNil(t, list[0].ProductID)
	require.Equal(t, productID, *list[0].ProductID)
}
