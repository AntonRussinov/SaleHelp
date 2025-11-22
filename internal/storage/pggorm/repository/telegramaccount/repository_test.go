package telegramaccount_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	tgrepo "salehelp/internal/storage/pggorm/repository/telegramaccount"
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

func setupRepo(t *testing.T) (*tgrepo.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := dbmock.MockExpectGORM(t)
	repo := tgrepo.NewRepository(slogdiscard.NewDiscardLogger(), db)

	return repo, mock
}

func JoinedColumns() []string {
	return model.TelegramAccountColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()

	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) TelegramAccount() domain.TelegramAccount {
	var a domain.TelegramAccount

	err := b.faker.Struct(&a)
	require.NoError(b.t, err)

	return a
}

func (b *DataBuilder) TelegramAccountRowValues(a domain.TelegramAccount) []driver.Value {
	return []driver.Value{
		a.ID,
		a.UserID,
		a.TelegramID,
		a.ChatID,
		a.Username,
		a.CreatedAt,
	}
}
