package merchant_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	merchantrepo "salehelp/internal/storage/pggorm/repository/merchant"
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

func setupRepo(t *testing.T) (*merchantrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := merchantrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	columns := []string{}
	columns = append(columns, model.MerchantColumns()...)

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

func (b *DataBuilder) Merchant() domain.Merchant {
	var merchant domain.Merchant

	err := b.faker.Struct(&merchant)

	require.NoError(b.t, err)

	return merchant
}

func (b *DataBuilder) MerchantRowValues(merch domain.Merchant) []driver.Value {
	return []driver.Value{
		merch.MerchantUID,
		merch.CompanyName,
		merch.Login,
		merch.Password,
		merch.IsEnabled,
		b.faker.Date(),
		b.faker.Date(),
		nil,
	}
}
