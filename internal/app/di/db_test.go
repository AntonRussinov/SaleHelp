package di_test

import (
	"context"
	"salehelp/internal/app/di"
	"salehelp/internal/config"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitDB(t *testing.T) {
	t.Parallel()

	t.Run("returns error on invalid db config", func(t *testing.T) {
		t.Parallel()

		prov := di.NewProvider(slogdiscard.NewDiscardLogger(), &config.ServiceConfig{})

		err := prov.InitDB(context.Background())
		require.Error(t, err)
	})
}
