package di_test

import (
	"context"
	"salehelp/internal/app/di"
	"salehelp/internal/config"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitCronRunner(t *testing.T) {
	t.Parallel()

	prov := di.NewProvider(slogdiscard.NewDiscardLogger(), &config.ServiceConfig{})

	err := prov.InitCronRunner(context.Background())
	require.NoError(t, err)
}
