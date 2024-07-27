package di_test

import (
	"reflect"
	"salehelp/internal/app/di"
	"salehelp/internal/config"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitPresenters(t *testing.T) {
	t.Parallel()

	prov := di.NewProvider(slogdiscard.NewDiscardLogger(), &config.ServiceConfig{})
	prov.InitPresenters()

	v := reflect.ValueOf(prov.Presenters)
	for i := range v.NumField() {
		assert.False(t, v.Field(i).IsNil())
	}
}
