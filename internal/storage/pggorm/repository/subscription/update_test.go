package subscription_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRepositoryExpireActiveByUserID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo, mock := setupSubscriptionRepo(t)

	userID := 42

	mock.ExpectExec(`UPDATE subscriptions SET status='expired'(.+)`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.ExpireActiveByUserID(ctx, userID)
	require.NoError(t, err)
}
