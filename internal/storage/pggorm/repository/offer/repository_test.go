package offer_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	cityPricerepo "salehelp/internal/storage/pggorm/repository/offer"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"salehelp/test/dbmock"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type DataBuilder struct {
	t     *testing.T
	faker *gofakeit.Faker
}

func setupRepo(t *testing.T) (*cityPricerepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := cityPricerepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	columns := []string{}
	columns = append(columns, model.OfferColumns()...)

	return columns
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	faker := gofakeit.New(1)

	return &DataBuilder{
		t:     t,
		faker: faker,
	}
}

func (b *DataBuilder) Offer() domain.Offer {
	var cityPrice domain.Offer

	err := b.faker.Struct(&cityPrice)

	require.NoError(b.t, err)

	return cityPrice
}

func (b *DataBuilder) OfferRowValues(cityPrice domain.Offer) []driver.Value {
	return []driver.Value{
		cityPrice.ID,
		cityPrice.OfferExternalID,
		cityPrice.SKU,
		cityPrice.ProductID,
		cityPrice.MerchantUID,
		cityPrice.ShopLink,
		cityPrice.FileID,
		cityPrice.OfferSource,
		cityPrice.ReindexReason,
		cityPrice.MatchesLastBatch,
		cityPrice.SkipDistinctCheck,
		cityPrice.Priority,
		cityPrice.Price,
		cityPrice.OldPrice,
		cityPrice.Available,
		cityPrice.MinPrice,
		cityPrice.MaxPrice,
		cityPrice.ArchivedAt,
		cityPrice.Modification,
		cityPrice.ColourCode,
		cityPrice.EditionDeadline,
		cityPrice.AnyPickup,
		cityPrice.AnyKaspiDelivery,
		cityPrice.AnyKaspiDeliveryLocal,
		cityPrice.AnyKaspiDeliveryExpress,
		cityPrice.AnyMerchantDelivery,
		cityPrice.BrandCode,
		cityPrice.SelectedCity,
		cityPrice.OfferError,
		b.faker.Date(),
		b.faker.Date(),
		nil,
	}

}
