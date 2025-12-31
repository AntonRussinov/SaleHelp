package productsnapshot_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	snaprepo "salehelp/internal/storage/pggorm/repository/productsnapshot"
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

func setupRepo(t *testing.T) (*snaprepo.Repository, sqlmock.Sqlmock) {
	db, mock := dbmock.MockExpectGORM(t)
	return snaprepo.NewRepository(slogdiscard.NewDiscardLogger(), db), mock
}

func JoinedColumns() []string {
	return model.ProductSnapshotColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) Snapshot() domain.ProductSnapshot {
	var s domain.ProductSnapshot
	err := b.faker.Struct(&s)
	require.NoError(b.t, err)
	return s
}

func (b *DataBuilder) ProductSnapshotRowValues(s domain.ProductSnapshot) []driver.Value {
	return []driver.Value{
		s.ID,
		s.ProductID,
		s.SnapshotTime,
		s.CurrencyCode,
		s.BasePrice,
		s.FinalPrice,
		s.ClubPrice,
		s.Discount,
		s.ClubDiscount,
		s.OldPrice,
		s.StockTotal,
		s.StockInWay,
		s.Rating,
		s.ReviewsCount,
		s.Sales24h,
		s.Sales7d,
		s.PositionSearch,
		s.IsBadTurnover,
		s.EditableSizePrice,
		s.Extra,
		s.CreatedAt,
	}
}
