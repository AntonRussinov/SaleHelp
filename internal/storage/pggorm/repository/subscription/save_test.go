package subscription_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
	}{
		{name: "executes insert query"},
		{name: "returns error on sql insert error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupSubscriptionRepo(t)
			sub := builder.Subscription()

			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "subscriptions"(.+)`).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			err := repo.Create(ctx, &sub)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
