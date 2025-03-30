package image_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/test/dbmock"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageRepositorySave(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		image   domain.Image
		sqlErr  error
		wantErr bool
	}{
		{name: "executes insert query", image: builder.Image()},
		{name: "returns error on sql insert error", image: builder.Image(), sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "images"(.+)`).WithArgs(saveImageSQLArgs(tt.image)...).
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			imageParam := tt.image

			err := repo.Save(ctx, &imageParam)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveImageSQLArgs(image domain.Image) []driver.Value {
	values := []driver.Value{
		image.ID,
		image.OfferID,
		image.ImagePath,
		dbmock.MatchTime(),
		dbmock.MatchTime(),
		nil,
	}

	return values
}
