package wbwarehouse_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	wbwarehouse "salehelp/internal/storage/pggorm/repository/wbwarehouse"
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

func setupRepo(t *testing.T) (*wbwarehouse.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := wbwarehouse.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	return model.WBWarehouseColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) Warehouse() domain.WBWarehouse {
	var w domain.WBWarehouse

	err := b.faker.Struct(&w)
	require.NoError(b.t, err)

	// Важно: WBID и Name у нас not null
	if w.WBID == 0 {
		w.WBID = b.faker.Int64()
	}
	if w.Name == "" {
		w.Name = b.faker.Company()
	}

	return w
}

func (b *DataBuilder) WarehouseRowValues(w domain.WBWarehouse) []driver.Value {
	return []driver.Value{
		w.ID,
		w.WBID,
		w.Name,
		w.Address,
		w.City,
		w.Longitude,
		w.Latitude,
		w.CargoType,
		w.DeliveryType,
		w.FederalDistrict,
		w.Selected,
		w.CreatedAt,
		w.UpdatedAt,
	}
}
