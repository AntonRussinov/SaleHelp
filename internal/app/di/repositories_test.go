package di_test

import (
	"reflect"
	"salehelp/internal/app/di"
	"salehelp/internal/config"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"salehelp/test/dbmock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRepositories(t *testing.T) {
	t.Parallel()

	prov := di.NewProvider(slogdiscard.NewDiscardLogger(), &config.ServiceConfig{})

	gormDB, _ := dbmock.MockGORM(t)
	db, err := gormDB.DB()

	require.NoError(t, err)

	prov.Database.DB = db
	prov.Database.GormDB = gormDB

	prov.InitRepositories()

	v := reflect.ValueOf(prov.Repositories)
	for i := range v.NumField() {
		assert.False(t, v.Field(i).IsNil())
	}
}
