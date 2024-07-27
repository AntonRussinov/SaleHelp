package azure

import (
	"context"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Cache struct {
	url string
	c   *jwk.Cache
}

const discoveryURLTemplate = "https://login.microsoftonline.com/%s/discovery/v2.0/keys?appid=%s"

//nolint:ireturn
func (cache *Cache) Get(ctx context.Context) (jwk.Set, error) {
	return cache.c.Get(ctx, cache.url)
}

func NewCache(tennantID, appID string) (*Cache, error) {
	cache := jwk.NewCache(context.Background())

	url := fmt.Sprintf(discoveryURLTemplate, tennantID, appID)

	err := cache.Register(url)
	if err != nil {
		return nil, err
	}

	_, err = cache.Refresh(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return &Cache{
		url: url,
		c:   cache,
	}, nil
}
