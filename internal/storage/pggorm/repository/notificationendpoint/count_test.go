package notificationendpoint_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCountEnabledByShopAndChannel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo, mock := setupRepo(t)

	shopID := 10
	channel := "telegram"

	mock.ExpectQuery(`SELECT count\(\*\) FROM "shop_notification_endpoints" WHERE shop_id = \$1 AND channel = \$2 AND enabled = TRUE`).
		WithArgs(shopID, channel).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	got, err := repo.CountEnabledByShopAndChannel(ctx, shopID, channel)
	require.NoError(t, err)
	require.Equal(t, 2, got)

	require.NoError(t, mock.ExpectationsWereMet())
}
