package alerttriggereddeliveries_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	deliveryrepo "salehelp/internal/storage/pggorm/repository/alerttriggereddeliveries"
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

func setupRepo(t *testing.T) (*deliveryrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := deliveryrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	return model.AlertTriggeredDeliveryColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) Delivery() domain.AlertTriggeredDelivery {
	var d domain.AlertTriggeredDelivery

	err := b.faker.Struct(&d)
	require.NoError(b.t, err)

	// NOT NULL
	if d.TriggeredID == 0 {
		d.TriggeredID = b.faker.Number(1, 100000)
	}
	if d.EndpointID == 0 {
		d.EndpointID = b.faker.Number(1, 100000)
	}
	if d.Status == "" {
		d.Status = "new"
	}

	// Ограничим даты, чтобы не было "нулей"
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now

	return d
}

func (b *DataBuilder) DeliveryRowValues(d domain.AlertTriggeredDelivery) []driver.Value {
	return []driver.Value{
		d.ID,
		d.TriggeredID,
		d.EndpointID,
		d.Status,
		d.SentAt,
		d.Attempts,
		d.LastError,
		d.LeaseUntil,
		d.LeasedBy,
		d.CreatedAt,
		d.UpdatedAt,
	}
}
