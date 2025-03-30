package image_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	imagerepo "salehelp/internal/storage/pggorm/repository/image"
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

func setupRepo(t *testing.T) (*imagerepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := imagerepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	columns := []string{}
	columns = append(columns, model.ImageColumns()...)

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

func (b *DataBuilder) Image() domain.Image {
	var image domain.Image

	err := b.faker.Struct(&image)

	require.NoError(b.t, err)

	return image
}

func (b *DataBuilder) ImageRowValues(image domain.Image) []driver.Value {
	return []driver.Value{
		image.ID,
		image.OfferID,
		image.ImagePath,
		b.faker.Date(),
		b.faker.Date(),
		nil,
	}

}
