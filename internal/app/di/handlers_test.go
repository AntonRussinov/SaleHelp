package di_test

import (
	"reflect"
	"salehelp/internal/app/di"
	"salehelp/internal/config"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitHandlers(t *testing.T) {
	t.Parallel()

	prov := di.NewProvider(slogdiscard.NewDiscardLogger(), &config.ServiceConfig{})
	err := prov.InitHandlers()
	require.NoError(t, err)

	v := reflect.ValueOf(prov.Handlers)
	for i := range v.NumField() {
		assert.False(t, v.Field(i).IsNil())
	}
}
