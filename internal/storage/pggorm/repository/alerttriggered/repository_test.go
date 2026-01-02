package alerttriggered_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	atrepo "salehelp/internal/storage/pggorm/repository/alerttriggered"
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

func setupRepo(t *testing.T) (*atrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := atrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	return model.AlertTriggeredColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) AlertTriggered() domain.AlertTriggered {
	var a domain.AlertTriggered

	err := b.faker.Struct(&a)
	require.NoError(b.t, err)

	return a
}

func (b *DataBuilder) AlertTriggeredRowValues(a domain.AlertTriggered) []driver.Value {
	return []driver.Value{
		a.ID,
		a.AlertID,
		a.ProductID,
		a.SnapshotID,
		a.ShopID,
		a.Message,
		a.ProductSizeID,
		a.AlertType,
		a.CreatedAt,
	}
}
