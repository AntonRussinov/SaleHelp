package user_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepositoryGetByID(t *testing.T) {
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			u := builder.User()
			u.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tc.found {
				rows.AddRow(builder.UserRowValues(u)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE id = \$1(.+)`).
				WithArgs(u.ID).
				WillReturnRows(rows).
				WillReturnError(tc.sqlErr)

			got, err := repo.GetByID(ctx, u.ID)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if !tc.found {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, u.ID, got.ID)
			require.Equal(t, u.Email, got.Email)
		})
	}
}

func TestUserRepositoryGetByEmail(t *testing.T) {
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			u := builder.User()

			rows := sqlmock.NewRows(JoinedColumns())
			if tc.found {
				rows.AddRow(builder.UserRowValues(u)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE email = \$1(.+)`).
				WithArgs(u.Email).
				WillReturnRows(rows).
				WillReturnError(tc.sqlErr)

			got, err := repo.GetByEmail(ctx, *u.Email)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if !tc.found {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, u.Email, got.Email)
		})
	}
}
