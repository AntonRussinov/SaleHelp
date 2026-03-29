// internal/storage/pggorm/repository/shop/repository_test.go
package shop_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	shoprepo "salehelp/internal/storage/pggorm/repository/shop"
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

func setupRepo(t *testing.T) (*shoprepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := shoprepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	columns := []string{}
	columns = append(columns, model.ShopColumns()...)

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

func (b *DataBuilder) Shop() domain.Shop {
	var shop domain.Shop

	err := b.faker.Struct(&shop)
	require.NoError(b.t, err)

	return shop
}

// порядок значений должен совпадать с model.ShopColumns()
func (b *DataBuilder) ShopRowValues(shop domain.Shop) []driver.Value {
	return []driver.Value{
		shop.ID,
		shop.UserID,
		shop.Marketplace,
		shop.ExternalID,
		shop.Name,
		shop.APIToken,
		shop.IsActive,
		shop.WBReportTimezone,
		shop.LastSyncAt,
		shop.FirstSyncNotifiedAt,
		shop.LastSyncStatusCode,
		shop.LastSyncError,
		shop.LastSyncFailedAt,
		shop.LastSyncOkAt,
		shop.CreatedAt,
		shop.UpdatedAt,
	}
}
