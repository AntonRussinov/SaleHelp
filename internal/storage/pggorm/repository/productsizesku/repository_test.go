package productsizesku_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	repo "salehelp/internal/storage/pggorm/repository/productsizesku"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"salehelp/test/dbmock"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type DataBuilder struct {
	t     *testing.T
	faker *gofakeit.Faker
}

func setupRepo(t *testing.T) (*repo.Repository, sqlmock.Sqlmock) {
	db, mock := dbmock.MockExpectGORM(t)
	return repo.NewRepository(slogdiscard.NewDiscardLogger(), db), mock
}

func JoinedColumns() []string {
	return model.ProductSizeSKUColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	return &DataBuilder{t: t, faker: gofakeit.New(1)}
}

func (b *DataBuilder) SKU() domain.ProductSizeSKU {
	var s domain.ProductSizeSKU
	err := b.faker.Struct(&s)
	require.NoError(b.t, err)

	if s.ProductSizeID == 0 {
		s.ProductSizeID = 1
	}
	if s.SKU == "" {
		s.SKU = b.faker.LetterN(8)
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}

	return s
}

func (b *DataBuilder) RowValues(s domain.ProductSizeSKU) []driver.Value {
	return []driver.Value{
		s.ID,
		s.ProductSizeID,
		s.SKU,
		s.CreatedAt,
	}
}
