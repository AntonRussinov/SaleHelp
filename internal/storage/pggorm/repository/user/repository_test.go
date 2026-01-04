package user_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	userrepo "salehelp/internal/storage/pggorm/repository/user"
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

func setupRepo(t *testing.T) (*userrepo.Repository, sqlmock.Sqlmock) {
	db, mock := dbmock.MockExpectGORM(t)
	return userrepo.NewRepository(slogdiscard.NewDiscardLogger(), db), mock
}

func JoinedColumns() []string {
	return model.UserColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) User() domain.User {
	var u domain.User
	err := b.faker.Struct(&u)
	require.NoError(b.t, err)
	return u
}

func (b *DataBuilder) UserRowValues(u domain.User) []driver.Value {
	return []driver.Value{
		u.ID,
		u.Email,
		u.PasswordHash,
		u.TelegramChatID,
		u.CreatedAt,
		u.UpdatedAt,
	}
}
