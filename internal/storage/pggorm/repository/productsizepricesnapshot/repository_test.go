package productsizepricesnapshot_test

import (
	"database/sql/driver"
	"testing"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	repo "salehelp/internal/storage/pggorm/repository/productsizepricesnapshot"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"salehelp/test/dbmock"

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
	return model.ProductSizePriceSnapshotColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) PriceSnapshot() domain.ProductSizePriceSnapshot {
	var s domain.ProductSizePriceSnapshot
	err := b.faker.Struct(&s)
	require.NoError(b.t, err)

	// важное: not null поля
	if s.ProductSnapshotID == 0 {
		s.ProductSnapshotID = 1
	}
	if s.ProductSizeID == 0 {
		s.ProductSizeID = 1
	}

	return s
}

func (b *DataBuilder) RowValues(s domain.ProductSizePriceSnapshot) []driver.Value {
	return []driver.Value{
		s.ID,
		s.ProductSnapshotID,
		s.ProductSizeID,
		s.BasePrice,
		s.FinalPrice,
		s.ClubPrice,
		s.Discount,
		s.ClubDiscount,
		s.Extra,
	}
}

func saveArgs(s domain.ProductSizePriceSnapshot) []driver.Value {
	return []driver.Value{
		s.ProductSnapshotID,
		s.ProductSizeID,
		s.BasePrice,
		s.FinalPrice,
		s.ClubPrice,
		s.Discount,
		s.ClubDiscount,
		s.Extra,
	}
}
