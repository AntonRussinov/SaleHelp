// internal/storage/pggorm/repository/shop/save_test.go
package shop_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShopRepositoryCreate(t *testing.T) {
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

			shop := builder.Shop()

			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "shops"(.+)`).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			shopParam := shop

			err := repo.Create(ctx, &shopParam)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveShopSQLArgs(shop domain.Shop) []driver.Value {
	// предполагаем, что id — BIGSERIAL и GORM его не передаёт в INSERT
	// если хочешь явно вставлять ID — добавь shop.ID первым
	return []driver.Value{
		shop.UserID,
		shop.Marketplace,
		shop.ExternalID,
		shop.Name,
		shop.APIToken,
		shop.IsActive,
		shop.LastSyncAt,
		shop.FirstSyncNotifiedAt,
		shop.LastSyncStatusCode,
		shop.LastSyncError,
		shop.LastSyncFailedAt,
		shop.LastSyncOkAt,
		shop.CreatedAt,
		shop.UpdatedAt,
	}
}
