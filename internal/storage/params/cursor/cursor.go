package cursor

import (
	"context"

	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

type (
	key string

	Params struct {
		Cursor *paginator.Cursor
		Limit  uint32
	}
)

const paramsKey key = "cursorParams"

func GetFromContext(ctx context.Context) *Params {
	params, ok := ctx.Value(paramsKey).(*Params)
	if !ok {
		return nil
	}

	return params
}

func SetToContext(ctx context.Context, params *Params) context.Context {
	return context.WithValue(ctx, paramsKey, params)
}
