package productsizepricesnapshot_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductSizePriceSnapshotRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
	}{
		{"success", nil, false},
		{"sql error", assert.AnError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository, mock := setupRepo(t)
			s := builder.PriceSnapshot()

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "product_size_price_snapshots"`).
				WithArgs(saveArgs(s)...).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			param := s
			err := repository.Create(ctx, &param)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
