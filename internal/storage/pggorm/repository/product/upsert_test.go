package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductRepositoryUpsert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	t.Run("creates new product when not exists", func(t *testing.T) {
		t.Parallel()

		repo, mock := setupRepo(t)

		p := builder.Product()
		p.ExternalID = "ext-123"

		// 1. SELECT returns no rows
		mock.ExpectQuery(`SELECT (.+) FROM "products"`).
			WithArgs(p.ShopID, p.ExternalID).
			WillReturnRows(sqlmock.NewRows([]string{})) // nil → no existing

		// 2. INSERT for new product
		args := saveProductSQLArgs(p)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "products"`).
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()

		err := repo.Upsert(ctx, &p)
		require.NoError(t, err)
	})

	t.Run("updates existing product", func(t *testing.T) {
		t.Parallel()

		repo, mock := setupRepo(t)

		existing := builder.Product()
		existing.ID = 42
		existing.ExternalID = "ext-222"

		updated := existing
		updated.Title = "Updated name"

		// 1. SELECT returns existing row
		mock.ExpectQuery(`SELECT (.+) FROM "products"`).
			WithArgs(existing.ShopID, existing.ExternalID).
			WillReturnRows(
				sqlmock.NewRows([]string{
					"id", "shop_id", "external_id", "title", "url",
					"image_url", "brand", "category", "description",
					"created_at", "updated_at",
				}).AddRow(
					existing.ID,
					existing.ShopID,
					existing.ExternalID,
					existing.Title,
					existing.URL,
					existing.ImageURL,
					existing.Brand,
					existing.Category,
					existing.Description,
					existing.CreatedAt,
					existing.UpdatedAt,
				),
			)

		// 2. UPDATE (Save)
		saveArgs := saveProductSQLArgs(updated)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "products" SET`).
			WithArgs(
				saveArgs[0], // shop_id
				saveArgs[1], // external_id
				saveArgs[2], // title
				saveArgs[3], // url
				saveArgs[4], // image_url
				saveArgs[5], // brand
				saveArgs[6], // category
				saveArgs[7], // created_at
				saveArgs[8], // updated_at
				existing.ID, // WHERE id = ?
			).
			WillReturnResult(sqlmock.NewResult(42, 1))
		mock.ExpectCommit()

		err := repo.Upsert(ctx, &updated)
		require.NoError(t, err)
	})

	t.Run("returns error when select returns error", func(t *testing.T) {
		t.Parallel()

		repo, mock := setupRepo(t)

		p := builder.Product()

		mock.ExpectQuery(`SELECT (.+) FROM "products"`).
			WithArgs(p.ShopID, p.ExternalID).
			WillReturnError(assert.AnError)

		err := repo.Upsert(ctx, &p)
		require.Error(t, err)
	})

	t.Run("returns error when insert fails", func(t *testing.T) {
		t.Parallel()

		repo, mock := setupRepo(t)

		p := builder.Product()

		mock.ExpectQuery(`SELECT (.+) FROM "products"`).
			WithArgs(p.ShopID, p.ExternalID).
			WillReturnRows(sqlmock.NewRows([]string{})) // not exists

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "products"`).
			WithArgs(saveProductSQLArgs(p)...).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err := repo.Upsert(ctx, &p)
		require.Error(t, err)
	})

	t.Run("returns error when update fails", func(t *testing.T) {
		t.Parallel()

		repo, mock := setupRepo(t)

		existing := builder.Product()
		existing.ID = 5

		mock.ExpectQuery(`SELECT (.+) FROM "products"`).
			WithArgs(existing.ShopID, existing.ExternalID).
			WillReturnRows(
				sqlmock.NewRows([]string{
					"id", "shop_id", "external_id", "title", "url",
					"image_url", "brand", "category", "description",
					"created_at", "updated_at",
				}).AddRow(
					existing.ID,
					existing.ShopID,
					existing.ExternalID,
					existing.Title,
					existing.URL,
					existing.ImageURL,
					existing.Brand,
					existing.Category,
					existing.Description,
					time.Now(), time.Now(),
				),
			)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "products"`).
			WithArgs(saveProductSQLArgs(existing)...).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err := repo.Upsert(ctx, &existing)
		require.Error(t, err)
	})
}
