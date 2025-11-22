package alert_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	alertrepo "salehelp/internal/storage/pggorm/repository/alert"
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

func setupRepo(t *testing.T) (*alertrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := alertrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	return model.AlertColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) Alert() domain.Alert {
	var a domain.Alert

	err := b.faker.Struct(&a)
	require.NoError(b.t, err)

	return a
}

func (b *DataBuilder) AlertRowValues(a domain.Alert) []driver.Value {
	return []driver.Value{
		a.ID,
		a.ShopID,
		a.ProductID,
		a.Type,
		a.Severity,
		a.Status,
		a.Message,
		a.Payload,
		a.TriggeredAt,
		a.ResolvedAt,
		a.AIRecommendationID,
	}
}
