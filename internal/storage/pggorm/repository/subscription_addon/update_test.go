package subscriptionaddon_test

import (
	"context"
	"salehelp/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionAddonRepositoryUpsert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo, mock := setupAddonRepo(t)

	add := domain.SubscriptionAddon{
		SubscriptionID: 10,
		Code:           "telegram_extra_chat",
		Quantity:       2,
		AmountMinor:    300000,
		CurrencyCode:   "KZT",
		Meta:           []byte(`{}`),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "subscription_addons"(.+)ON CONFLICT \("subscription_id","code"\) DO UPDATE SET(.+)`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Upsert(ctx, &add)
	require.NoError(t, err)
	assert.True(t, add.ID >= 0) // если gorm вернёт id
}
