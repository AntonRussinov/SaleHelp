package subscription_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRepositoryExpireActiveByShopID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo, mock := setupSubscriptionRepo(t)

	shopID := 123

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "subscriptions" SET (.+) WHERE (.+)`).
		WithArgs("expired", sqlmock.AnyArg(), shopID, "active").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	affected, err := repo.ExpireActiveByShopID(ctx, shopID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)
}
