package productsizesnapshot_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	repo "salehelp/internal/storage/pggorm/repository/productsizesnapshot"
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

func setupRepo(t *testing.T) (*repo.Repository, sqlmock.Sqlmock) {
	db, mock := dbmock.MockExpectGORM(t)
	return repo.NewRepository(slogdiscard.NewDiscardLogger(), db), mock
}

func JoinedColumns() []string {
	return model.ProductSizeSnapshotColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) SizeSnapshot() domain.ProductSizeSnapshot {
	var s domain.ProductSizeSnapshot
	err := b.faker.Struct(&s)
	require.NoError(b.t, err)
	return s
}

func (b *DataBuilder) RowValues(s domain.ProductSizeSnapshot) []driver.Value {
	return []driver.Value{
		s.ID,
		s.ProductSnapshotID,
		s.ProductSizeID,
		s.SKU,
		s.Stock,
		s.StockByWarehouse,
		s.Extra,
	}
}
