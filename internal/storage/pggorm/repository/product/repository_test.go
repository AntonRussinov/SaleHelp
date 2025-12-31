package product_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	productrepo "salehelp/internal/storage/pggorm/repository/product"
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

func setupRepo(t *testing.T) (*productrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := productrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	columns := []string{}
	columns = append(columns, model.ProductColumns()...)

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

func (b *DataBuilder) Product() domain.Product {
	var product domain.Product

	err := b.faker.Struct(&product)
	require.NoError(b.t, err)

	return product
}

// порядок значений должен совпадать с model.ProductColumns()
func (b *DataBuilder) ProductRowValues(product domain.Product) []driver.Value {
	return []driver.Value{
		product.ID,
		product.ShopID,
		product.ExternalID,
		product.Title,
		product.URL,
		product.ImageURL,
		product.Brand,
		product.Category,
		product.Description,
		product.CreatedAt,
		product.UpdatedAt,
	}
}
