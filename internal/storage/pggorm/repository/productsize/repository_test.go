package productsize_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	repo "salehelp/internal/storage/pggorm/repository/productsize"
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

func JoinedColumns() []string { return model.ProductSizeColumns() }

func NewDataBuilder(t *testing.T) *DataBuilder {
	return &DataBuilder{t: t, faker: gofakeit.New(1)}
}

func (b *DataBuilder) ProductSize() domain.ProductSize {
	var s domain.ProductSize
	err := b.faker.Struct(&s)
	require.NoError(b.t, err)

	// pointers: faker не всегда заполняет *string как нужно
	if s.TechSizeName == nil {
		v := b.faker.LetterN(2)
		s.TechSizeName = &v
	}

	if s.ExternalID == "" {
		s.ExternalID = b.faker.DigitN(8)
	}
	if s.ProductID == 0 {
		s.ProductID = 1
	}

	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = time.Now()
	}

	return s
}

func (b *DataBuilder) RowValues(s domain.ProductSize) []driver.Value {
	return []driver.Value{
		s.ID,
		s.ProductID,
		s.ExternalID,
		s.TechSizeName,
		s.CreatedAt,
		s.UpdatedAt,
	}
}
