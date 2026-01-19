package productsize_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductSizeRepositoryUpsert(t *testing.T) {
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

			repo, mock := setupRepo(t)
			ps := builder.ProductSize()
			ps.ID = 0 // new

			// GORM upsert = INSERT ... ON CONFLICT ... DO UPDATE ... RETURNING ...
			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "product_sizes"`).
				WithArgs(
					sqlmock.AnyArg(), // created_at
					sqlmock.AnyArg(), // updated_at
					ps.ProductID,
					ps.ExternalID,
					ps.TechSizeName,
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).AddRow(10),
				).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			param := ps
			err := repo.Upsert(ctx, &param)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotZero(t, param.ID)
		})
	}
}
