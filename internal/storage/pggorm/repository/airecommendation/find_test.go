package airecommendation_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAIRecommendationRepositoryGetByID(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			rec := builder.Recommendation()
			rec.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.RowValues(rec)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "ai_recommendations" WHERE id = \$1(.+)`).
				WithArgs(rec.ID).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, rec.ID)

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
			require.Equal(t, rec.ID, got.ID)
			require.Equal(t, rec.ShopID, got.ShopID)
			require.Equal(t, rec.Type, got.Type)
		})
	}
}

func TestAIRecommendationRepositoryListByShop(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	shopID := 100
	r1 := builder.Recommendation()
	r1.ShopID = shopID
	r2 := builder.Recommendation()
	r2.ShopID = shopID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.RowValues(r1)...).
		AddRow(builder.RowValues(r2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "ai_recommendations" WHERE shop_id = \$1(.+)`).
		WithArgs(shopID).
		WillReturnRows(rows)

	list, err := repo.ListByShop(ctx, shopID)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, shopID, list[0].ShopID)
}

func TestAIRecommendationRepositoryListByProduct(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	productID := 777

	r1 := builder.Recommendation()
	r1.ProductID = &productID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.RowValues(r1)...)

	mock.ExpectQuery(`SELECT (.+) FROM "ai_recommendations" WHERE product_id = \$1(.+)`).
		WithArgs(productID).
		WillReturnRows(rows)

	list, err := repo.ListByProduct(ctx, productID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, productID, *list[0].ProductID)
}

func TestAIRecommendationRepositoryListByAlert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	alertID := 50

	r1 := builder.Recommendation()
	r1.AlertID = &alertID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.RowValues(r1)...)

	mock.ExpectQuery(`SELECT (.+) FROM "ai_recommendations" WHERE alert_id = \$1(.+)`).
		WithArgs(alertID).
		WillReturnRows(rows)

	list, err := repo.ListByAlert(ctx, alertID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, alertID, *list[0].AlertID)
}
