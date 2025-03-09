package pagination

import "context"

type (
	key string

	Params struct {
		Offset uint32
		Limit  uint32
	}
)

const paramsKey key = "paginationParams"

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

func HasNext(ctx context.Context, totalCount int64) bool {
	params := GetFromContext(ctx)
	if params == nil || params.Limit == 0 {
		return false
	}

	requested := params.Offset + params.Limit

	return int64(requested) < totalCount
}
