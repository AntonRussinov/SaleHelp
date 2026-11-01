package sellerwarehouse_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	sellerwarehouse "salehelp/internal/storage/pggorm/repository/sellerwarehouse"
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

func setupRepo(t *testing.T) (*sellerwarehouse.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := sellerwarehouse.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	return model.SellerWarehouseColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) Warehouse() domain.SellerWarehouse {
	var w domain.SellerWarehouse

	err := b.faker.Struct(&w)
	require.NoError(b.t, err)

	if w.ShopID == 0 {
		w.ShopID = b.faker.Number(1, 1000)
	}
	if w.WBOfficeID == 0 {
		w.WBOfficeID = b.faker.Int64()
	}
	if w.Name == "" {
		w.Name = b.faker.Company()
	}

	return w
}

func (b *DataBuilder) WarehouseRowValues(w domain.SellerWarehouse) []driver.Value {
	return []driver.Value{
		w.ID,
		w.ShopID,
		w.WBOfficeID,
		w.WBID,
		w.Name,
		w.CargoType,
		w.DeliveryType,
		w.IsDeleting,
		w.IsProcessing,
		w.CreatedAt,
		w.UpdatedAt,
	}
}
