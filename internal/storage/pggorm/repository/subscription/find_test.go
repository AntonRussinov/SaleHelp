package subscription_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRepositoryGetActiveByShopID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo, mock := setupSubscriptionRepo(t)

	now := time.Now()
	shopID := 123

	cols := []string{
		"id", "shop_id", "plan", "status", "valid_from", "valid_to",
		"period_months", "base_amount_minor", "currency_code",
		"paid_at", "source", "comment", "created_at", "updated_at",
	}

	mock.ExpectQuery(`SELECT (.+) FROM "subscriptions"(.+)`).
		WithArgs(shopID, "active", now).
		WillReturnRows(sqlmock.NewRows(cols).
			AddRow(
				1, shopID, "basic", "active",
				now.Add(-time.Hour), now.Add(time.Hour),
				1, int64(150000), "RUB",
				nil, "manual", nil, now, now,
			))

	got, err := repo.GetActiveByShopID(ctx, shopID, now)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, shopID, got.ShopID)
	assert.Equal(t, "active", got.Status)
}
