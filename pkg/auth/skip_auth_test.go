package auth_test

import (
	"salehelp/pkg/auth"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkipAuth(t *testing.T) {
	t.Parallel()

	authenticator := auth.NewSkipAuthenticator()
	user, _, err := authenticator.AuthenticateWithHeader(nil)
	require.NoError(t, err)
	require.NotNil(t, user)

	user, _, err = authenticator.AuthenticateWithQuery(nil)
	require.NoError(t, err)
	require.NotNil(t, user)
}
