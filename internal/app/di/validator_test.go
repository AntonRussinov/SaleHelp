package di_test

import (
	"salehelp/internal/app/di"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitValidator(t *testing.T) {
	t.Parallel()

	prov := &di.Provider{}

	err := prov.InitValidator()
	require.NoError(t, err)
	require.NotNil(t, prov.Validator)
}
