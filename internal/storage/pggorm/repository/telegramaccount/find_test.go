package telegramaccount_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegramAccountRepositoryGetByUserID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
		found   bool
		userID  int
	}{
		{name: "returns account", found: true, userID: 10},
		{name: "returns nil when not found", found: false, userID: 10},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true, userID: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			acc := builder.TelegramAccount()
			acc.UserID = tt.userID

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.TelegramAccountRowValues(acc)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "telegram_accounts" WHERE user_id = \$1(.+)`).
				WithArgs(tt.userID).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByUserID(ctx, tt.userID)

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
			require.Equal(t, acc.UserID, got.UserID)
			require.Equal(t, acc.TelegramID, got.TelegramID)
		})
	}
}

func TestTelegramAccountRepositoryGetByTelegramID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name       string
		sqlErr     error
		wantErr    bool
		found      bool
		telegramID int64
	}{
		{name: "returns account", found: true, telegramID: 123456789},
		{name: "returns nil when not found", found: false, telegramID: 123456789},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true, telegramID: 123456789},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			acc := builder.TelegramAccount()
			acc.TelegramID = tt.telegramID

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.TelegramAccountRowValues(acc)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "telegram_accounts" WHERE telegram_id = \$1(.+)`).
				WithArgs(tt.telegramID).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByTelegramID(ctx, tt.telegramID)

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
			require.Equal(t, acc.TelegramID, got.TelegramID)
			require.Equal(t, acc.UserID, got.UserID)
		})
	}
}
