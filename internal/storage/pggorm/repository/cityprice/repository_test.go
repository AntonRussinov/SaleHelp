package cityprice_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	cityPricerepo "salehelp/internal/storage/pggorm/repository/cityprice"
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
	columns = append(columns, model.CityPriceColumns()...)

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

func (b *DataBuilder) CityPrice() domain.CityPrice {
	var cityPrice domain.CityPrice

	err := b.faker.Struct(&cityPrice)

	require.NoError(b.t, err)

	return cityPrice
}

func (b *DataBuilder) CityPriceRowValues(cityPrice domain.CityPrice) []driver.Value {
	return []driver.Value{
		cityPrice.ID,
		cityPrice.MerchantUID,
		cityPrice.OfferID,
		cityPrice.CityID,
		cityPrice.Price,
		cityPrice.OldPrice,
		b.faker.Date(),
		b.faker.Date(),
		nil,
	}

}
