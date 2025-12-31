package productsizepricesnapshot_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestProductSizePriceSnapshotRepositoryDeleteBySnapshotID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repository, mock := setupRepo(t)

	snapshotID := 123

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "product_size_price_snapshots" WHERE product_snapshot_id = \$1`).
		WithArgs(snapshotID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	err := repository.DeleteBySnapshotID(ctx, snapshotID)
	require.NoError(t, err)
}
