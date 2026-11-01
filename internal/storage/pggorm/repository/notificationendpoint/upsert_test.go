package notificationendpoint_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationEndpointRepositoryUpsert(t *testing.T) {
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

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			e := builder.Endpoint()

			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "shop_notification_endpoints"`).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			param := e
			err := repo.Upsert(ctx, &param)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
