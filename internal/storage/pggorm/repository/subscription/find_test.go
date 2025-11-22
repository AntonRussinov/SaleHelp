package subscription_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRepositoryGetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
		found   bool
	}{
		{name: "returns subscription", found: true},
		{name: "returns nil when not found", found: false},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			sub := builder.Subscription()
			sub.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.SubscriptionRowValues(sub)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "subscriptions" WHERE id = \$1(.+)`).
				WithArgs(sub.ID).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, sub.ID)

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
			require.Equal(t, sub.ID, got.ID)
			require.Equal(t, sub.UserID, got.UserID)
			require.Equal(t, sub.Plan, got.Plan)
		})
	}
}

func TestSubscriptionRepositoryListByUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	userID := 42

	s1 := builder.Subscription()
	s1.UserID = userID

	s2 := builder.Subscription()
	s2.UserID = userID

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.SubscriptionRowValues(s1)...).
		AddRow(builder.SubscriptionRowValues(s2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "subscriptions" WHERE user_id = \$1(.+)`).
		WithArgs(userID).
		WillReturnRows(rows)

	list, err := repo.ListByUser(ctx, userID)

	require.NoError(t, err)
	require.Len(t, list, 2)

	require.Equal(t, userID, list[0].UserID)
	require.Equal(t, userID, list[1].UserID)
}
