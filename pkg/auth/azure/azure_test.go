package azure_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"salehelp/pkg/auth/azure"
	"testing"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	Keys struct {
		Private  jwk.Key
		Wrongkid jwk.Key
		Public   jwk.Key
	}

	KeySetCache struct {
		set       jwk.Set
		errOncall error
	}
)

func TestValidateToken(t *testing.T) {
	t.Parallel()

	keys := keys(t)

	emptyCache := keySetCache(t)
	errorCache := errorKeySetCache(t)
	cache := keySetCache(t)
	err := cache.set.AddKey(keys.Public)
	require.NoError(t, err)

	tests := []struct {
		name    string
		cache   *KeySetCache
		wantErr bool
		alg     jwa.KeyAlgorithm
		key     any
	}{
		{name: "empty cache", cache: emptyCache, wantErr: true, alg: jwa.RS256, key: keys.Private},
		{name: "error cache", cache: errorCache, wantErr: true, alg: jwa.RS256, key: keys.Private},
		{name: "wrong alg", cache: cache, wantErr: true, alg: jwa.HS256, key: []byte("6ZXMcKhTstHzvVtDuovPsvfjxIxfhkn0boHr88vl8jc=")},
		{name: "wrong kid", cache: cache, wantErr: true, alg: jwa.RS256, key: keys.Wrongkid},
		{name: "valid", cache: cache, wantErr: false, alg: jwa.RS256, key: keys.Private},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			auth := azure.NewAuthentificator(tt.cache)

			_, err := auth.ValidateToken(context.Background(), token(t, tt.alg, tt.key))
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func token(t *testing.T, alg jwa.KeyAlgorithm, key any) string {
	t.Helper()

	token, err := jwt.NewBuilder().JwtID("1").Build()
	require.NoError(t, err)

	signed, err := jwt.Sign(token, jwt.WithKey(alg, key))
	require.NoError(t, err)

	return string(signed)
}

func keys(t *testing.T) *Keys {
	t.Helper()

	rsakey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	private, err := jwk.FromRaw(rsakey)
	require.NoError(t, err)

	err = private.Set(jwk.KeyIDKey, "1")
	require.NoError(t, err)

	err = private.Set(jwk.AlgorithmKey, jwa.RS256)
	require.NoError(t, err)

	public, err := jwk.PublicKeyOf(private)
	require.NoError(t, err)

	wrongkid, err := jwk.FromRaw(rsakey)
	require.NoError(t, err)

	err = wrongkid.Set(jwk.KeyIDKey, "2")
	require.NoError(t, err)

	err = wrongkid.Set(jwk.AlgorithmKey, jwa.RS256)
	require.NoError(t, err)

	return &Keys{
		Private:  private,
		Public:   public,
		Wrongkid: wrongkid,
	}
}

func keySetCache(t *testing.T) *KeySetCache {
	t.Helper()

	set := jwk.NewSet()

	return &KeySetCache{
		set: set,
	}
}

func errorKeySetCache(t *testing.T) *KeySetCache {
	t.Helper()

	set := jwk.NewSet()

	return &KeySetCache{
		set:       set,
		errOncall: assert.AnError,
	}
}

func (ks *KeySetCache) Get(_ context.Context) (jwk.Set, error) {
	if ks.errOncall != nil {
		return nil, ks.errOncall
	}

	return ks.set, nil
}
