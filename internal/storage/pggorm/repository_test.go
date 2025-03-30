package pggorm_test

import (
	"salehelp/internal/storage"
	"salehelp/internal/storage/pggorm"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPgGormHandleError(t *testing.T) {
	t.Parallel()

	repo := &pggorm.Repository{}
	repo.Log = slogdiscard.NewDiscardLogger()
	fn := "caller"

	t.Run("returns storage not found error on gorm not found", func(t *testing.T) {
		t.Parallel()

		err := repo.HandleError(fn, gorm.ErrRecordNotFound)
		require.ErrorIs(t, err, storage.ErrNotFound)
	})

	t.Run("returns storage duplicate error on pg unique constrain error", func(t *testing.T) {
		t.Parallel()

		pgerr := &pgconn.PgError{
			Code:           "23505",
			ConstraintName: "idx_applications_name",
		}

		err := repo.HandleError(fn, pgerr)

		var target *storage.DuplicateError

		require.ErrorAs(t, err, &target)
		initialErr := target.Unwrap()
		require.Equal(t, pgerr, initialErr)
	})

	t.Run("returns wrapped initial pg error on unique constrain error if unknown constraint name", func(t *testing.T) {
		t.Parallel()

		pgerr := &pgconn.PgError{
			Code:           "23505",
			ConstraintName: "some_unknown_constraint",
		}

		err := repo.HandleError(fn, pgerr)

		var target *pgconn.PgError

		require.ErrorAs(t, err, &target)
		require.Equal(t, pgerr, target)
	})

	t.Run("returns storage foreign key error on pg unique constrain error", func(t *testing.T) {
		t.Parallel()

		pgerr := &pgconn.PgError{
			Code:           "23503",
			ConstraintName: "applications_function_group_id_fkey",
		}

		err := repo.HandleError(fn, pgerr)

		var target *storage.ForeignKeyError

		require.ErrorAs(t, err, &target)
		initialErr := target.Unwrap()
		require.Equal(t, pgerr, initialErr)
	})

	t.Run("returns wrapped initial pg error on foreign key error if unknown constraint name", func(t *testing.T) {
		t.Parallel()

		pgerr := &pgconn.PgError{
			Code:           "23503",
			ConstraintName: "some_unknown_constraint",
		}

		err := repo.HandleError(fn, pgerr)

		var target *pgconn.PgError

		require.ErrorAs(t, err, &target)
		require.Equal(t, pgerr, target)
	})

	t.Run("returns wrapped initial pg error on any other code", func(t *testing.T) {
		t.Parallel()

		pgerr := &pgconn.PgError{
			Code: "Some error code",
		}

		err := repo.HandleError(fn, pgerr)

		var target *pgconn.PgError

		require.ErrorAs(t, err, &target)
		require.Equal(t, pgerr, target)
	})
}
