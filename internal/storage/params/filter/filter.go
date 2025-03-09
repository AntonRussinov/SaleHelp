package filter

import "context"

type (
	key string

	Param struct {
		Field       string
		CompareType CompareType
		Value       string
		Values      []string
	}
)

const paramsKey key = "filterParams"

func GetFromContext(ctx context.Context) []Param {
	filters, ok := ctx.Value(paramsKey).([]Param)
	if !ok {
		return nil
	}

	return filters
}

func SetToContext(ctx context.Context, params []Param) context.Context {
	return context.WithValue(ctx, paramsKey, params)
}
