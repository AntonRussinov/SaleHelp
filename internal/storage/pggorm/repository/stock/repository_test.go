package stock_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	stockrepo "salehelp/internal/storage/pggorm/repository/stock"
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

func setupRepo(t *testing.T) (*stockrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := stockrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	columns := []string{}
	columns = append(columns, model.StockColumns()...)

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

func (b *DataBuilder) Stock() domain.Stock {
	var stock domain.Stock

	err := b.faker.Struct(&stock)

	require.NoError(b.t, err)

	return stock
}

func (b *DataBuilder) StockRowValues(stock domain.Stock) []driver.Value {
	return []driver.Value{
		stock.ID,
		stock.MerchantUID,
		stock.OfferID,
		stock.CityID,
		stock.StoreID,
		stock.Available,
		stock.PreOrder,
		stock.StockCount,
		stock.Enabled,
		b.faker.Date(),
		b.faker.Date(),
		nil,
	}
}
