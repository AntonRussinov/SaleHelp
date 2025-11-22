package cronjob_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	cronrepo "salehelp/internal/storage/pggorm/repository/cronjob"
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

func setupRepo(t *testing.T) (*cronrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := cronrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	return model.CronJobColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) CronJob() domain.CronJob {
	var c domain.CronJob

	err := b.faker.Struct(&c)
	require.NoError(b.t, err)

	return c
}

func (b *DataBuilder) CronJobRowValues(c domain.CronJob) []driver.Value {
	return []driver.Value{
		c.ID,
		c.Name,
		c.Status,
		c.ErrorMessage,
		c.StartedAt,
		c.FinishedAt,
	}
}
