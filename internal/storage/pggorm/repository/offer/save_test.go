package offer_test

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

func TestOfferRepositorySave(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		offer   domain.Offer
		sqlErr  error
		wantErr bool
	}{
		{name: "executes insert query", offer: builder.Offer()},
		{name: "returns error on sql insert error", offer: builder.Offer(), sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "offers"(.+)`).WithArgs(saveOfferSQLArgs(tt.offer)...).
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			offerParam := tt.offer

			err := repo.Save(ctx, &offerParam)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveOfferSQLArgs(offer domain.Offer) []driver.Value {
	values := []driver.Value{
		offer.ID,
		offer.OfferExternalID,
		offer.SKU,
		offer.ProductID,
		offer.MerchantUID,
		offer.ShopLink,
		offer.FileID,
		offer.OfferSource,
		offer.ReindexReason,
		offer.MatchesLastBatch,
		offer.SkipDistinctCheck,
		offer.Priority,
		offer.Price,
		offer.OldPrice,
		offer.Available,
		offer.MinPrice,
		offer.MaxPrice,
		offer.ArchivedAt,
		offer.Modification,
		offer.ColourCode,
		offer.EditionDeadline,
		offer.AnyPickup,
		offer.AnyKaspiDelivery,
		offer.AnyKaspiDeliveryLocal,
		offer.AnyKaspiDeliveryExpress,
		offer.AnyMerchantDelivery,
		offer.BrandCode,
		offer.SelectedCity,
		offer.OfferError,
		dbmock.MatchTime(),
		dbmock.MatchTime(),
		nil,
	}

	return values
}
