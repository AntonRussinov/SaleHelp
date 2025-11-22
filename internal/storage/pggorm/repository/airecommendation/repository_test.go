package airecommendation_test

import (
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	aiRepo "salehelp/internal/storage/pggorm/repository/airecommendation"
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

func setupRepo(t *testing.T) (*aiRepo.Repository, sqlmock.Sqlmock) {
	db, mock := dbmock.MockExpectGORM(t)
	return aiRepo.NewRepository(slogdiscard.NewDiscardLogger(), db), mock
}

func JoinedColumns() []string {
	return model.AIRecommendationColumns()
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	return &DataBuilder{
		t:     t,
		faker: gofakeit.New(1),
	}
}

func (b *DataBuilder) Recommendation() domain.AIRecommendation {
	var r domain.AIRecommendation
	err := b.faker.Struct(&r)
	require.NoError(b.t, err)
	return r
}

func (b *DataBuilder) RowValues(r domain.AIRecommendation) []driver.Value {
	return []driver.Value{
		r.ID,
		r.ShopID,
		r.ProductID,
		r.AlertID,
		r.Type,
		r.Model,
		r.Prompt,
		r.Content,
		r.Meta,
		r.CreatedAt,
	}
}
