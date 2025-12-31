package notificationendpoint_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	endpointrepo "salehelp/internal/storage/pggorm/repository/notificationendpoint"
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

func setupRepo(t *testing.T) (*endpointrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := endpointrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	return model.ShopNotificationEndpointColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) Endpoint() domain.NotificationEndpoint {
	var e domain.NotificationEndpoint
	err := b.faker.Struct(&e)
	require.NoError(b.t, err)

	// NOT NULL поля
	if e.ShopID == 0 {
		e.ShopID = b.faker.Number(1, 1000)
	}
	if e.Channel == "" {
		e.Channel = "telegram"
	}
	if e.Address == "" {
		e.Address = b.faker.UUID()
	}

	return e
}

func (b *DataBuilder) EndpointRowValues(e domain.NotificationEndpoint) []driver.Value {
	return []driver.Value{
		e.ID,
		e.ShopID,
		e.Channel,
		e.Address,
		e.Enabled,
		e.Label,
		e.CreatedAt,
		e.UpdatedAt,
	}
}
