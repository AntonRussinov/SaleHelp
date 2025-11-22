package subscription_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRepositoryUpdate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
	}{
		{name: "executes update query"},
		{name: "returns error on sql update error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			sub := builder.Subscription()
			sub.ID = 1 // чтобы GORM точно делал UPDATE по PK

			mock.ExpectBegin()
			mock.ExpectExec(`UPDATE "subscriptions" SET`).
				// аргументы не проверяем, чтобы тест не зависел от внутренностей GORM
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			subParam := sub
			err := repo.Update(ctx, &subParam)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
