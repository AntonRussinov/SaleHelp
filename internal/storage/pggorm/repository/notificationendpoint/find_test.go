package notificationendpoint_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationEndpointRepositoryGetByID(t *testing.T) {
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

			e := builder.Endpoint()
			e.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.EndpointRowValues(e)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "shop_notification_endpoints" WHERE id = \$1`).
				WithArgs(e.ID, 1).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, e.ID)

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
			require.Equal(t, e.ID, got.ID)
			require.Equal(t, e.ShopID, got.ShopID)
		})
	}
}

func TestNotificationEndpointRepositoryListByShopID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	e1 := builder.Endpoint()
	e2 := builder.Endpoint()
	shopID := e1.ShopID
	e2.ShopID = shopID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.EndpointRowValues(e1)...).
		AddRow(builder.EndpointRowValues(e2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "shop_notification_endpoints" WHERE shop_id = \$1`).
		WithArgs(shopID).
		WillReturnRows(rows)

	list, err := repo.ListByShopID(ctx, shopID, false)

	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestNotificationEndpointRepositoryListByShopID_OnlyEnabled(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	e1 := builder.Endpoint()
	e1.Enabled = true
	shopID := e1.ShopID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.EndpointRowValues(e1)...)

	mock.ExpectQuery(`SELECT (.+) FROM "shop_notification_endpoints" WHERE shop_id = \$1 AND enabled = TRUE`).
		WithArgs(shopID).
		WillReturnRows(rows)

	list, err := repo.ListByShopID(ctx, shopID, true)

	require.NoError(t, err)
	require.Len(t, list, 1)
	require.True(t, list[0].Enabled)
}
