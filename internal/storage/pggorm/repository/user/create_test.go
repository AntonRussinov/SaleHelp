package user_test

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
	}{
		{name: "success"},
		{name: "sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)
			u := builder.User()

			args := []driver.Value{
				u.Email,
				u.PasswordHash,
				u.CreatedAt,
				u.UpdatedAt,
			}

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "users"`).
				WithArgs(args...).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tc.sqlErr)

			if tc.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			uParam := u
			err := repo.Create(ctx, &uParam)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
