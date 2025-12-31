package productsizesnapshot_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestProductSizeSnapshotRepositoryDeleteBySnapshotID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo, mock := setupRepo(t)

	snapshotID := 123

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "product_size_snapshots" WHERE product_snapshot_id = \$1`).
		WithArgs(snapshotID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	err := repo.DeleteBySnapshotID(ctx, snapshotID)
	require.NoError(t, err)
}
