package azure

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type (
	Config struct {
		AppID     string
		TennantID string
	}

	KeySetCache interface {
		Get(ctx context.Context) (jwk.Set, error)
	}

	Authentificator struct {
		KeySetCache KeySetCache
	}
)

func NewAuthentificator(cache KeySetCache) *Authentificator {
	return &Authentificator{
		KeySetCache: cache,
	}
}

func (a *Authentificator) ValidateToken(ctx context.Context, tokenString string) (map[string]any, error) {
	keySet, err := a.KeySetCache.Get(ctx)
	if err != nil {
		return nil, InvalidTokenError("failed to get key set")
	}

	token, err := jwt.ParseString(tokenString, jwt.WithKeySet(keySet, jws.WithInferAlgorithmFromKey(true)))
	if err != nil {
		return nil, InvalidTokenError("could not validate token")
	}

	claims, err := token.AsMap(ctx)
	if err != nil {
		return nil, InvalidTokenError("could not parse claims")
	}

	return claims, nil
}
