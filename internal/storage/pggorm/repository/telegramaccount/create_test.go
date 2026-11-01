package telegramaccount_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegramAccountRepositoryCreate(t *testing.T) {
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

			repo, mock := setupRepo(t)
			acc := builder.TelegramAccount()

			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "telegram_accounts"(.+)`).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			accParam := acc

			err := repo.Create(ctx, &accParam)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveTelegramAccountSQLArgs(a domain.TelegramAccount) []driver.Value {
	// предполагаем, что ID — BIGSERIAL и GORM его не вставляет
	return []driver.Value{
		a.UserID,
		a.TelegramID,
		a.ChatID,
		a.Username,
		a.CreatedAt,
	}
}
