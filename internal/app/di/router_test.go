package di_test

import (
	"salehelp/internal/app/di"
	"salehelp/internal/config"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitRouter(t *testing.T) {
	t.Parallel()

	prov := di.NewProvider(slogdiscard.NewDiscardLogger(), &config.ServiceConfig{})
	err := prov.InitHandlers()
	require.NoError(t, err)
	prov.InitRouter()

	require.NotNil(t, prov.Router)
}
