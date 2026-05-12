package subscription_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRepositoryGetActiveByUserID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo, mock := setupSubscriptionRepo(t)

	now := time.Now()
	userID := 42
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)

	cols := []string{
		"id", "user_id", "tier", "status", "period_months",
		"started_at", "expires_at", "trial_ends_at",
		"ai_quota", "ai_used", "ai_resets_at",
		"max_shops", "extra_shops", "price_minor",
		"created_at", "updated_at",
	}

	mock.ExpectQuery(`SELECT (.+) FROM "subscriptions"(.+)`).
		WithArgs(userID, sqlmock.AnyArg(), 1).
		WillReturnRows(sqlmock.NewRows(cols).
			AddRow(
				1, userID, "business", "active", 1,
				now.Add(-time.Hour), now.Add(30*24*time.Hour), nil,
				200, 5, nextMonth,
				2, 0, 149000,
				now, now,
			))

	got, err := repo.GetActiveByUserID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, userID, got.UserID)
	assert.Equal(t, "business", got.Tier)
	assert.Equal(t, "active", got.Status)
}
