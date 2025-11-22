package subscription_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	subrepo "salehelp/internal/storage/pggorm/repository/subscription"
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

func setupRepo(t *testing.T) (*subrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := subrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	return model.SubscriptionColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) Subscription() domain.Subscription {
	var s domain.Subscription

	err := b.faker.Struct(&s)
	require.NoError(b.t, err)

	return s
}

func (b *DataBuilder) SubscriptionRowValues(s domain.Subscription) []driver.Value {
	return []driver.Value{
		s.ID,
		s.UserID,
		s.Plan,
		s.Status,
		s.ValidFrom,
		s.ValidTo,
		s.ExternalCustomerID,
		s.ExternalSubscriptionID,
		s.CreatedAt,
		s.UpdatedAt,
	}
}
