package querystring

import "context"

type (
	key string
)

const paramKey key = "queryString"

func GetFromContext(ctx context.Context) string {
	qs, ok := ctx.Value(paramKey).(string)
	if !ok {
		return ""
	}

	return qs
}

func SetToContext(ctx context.Context, queryString string) context.Context {
	return context.WithValue(ctx, paramKey, queryString)
}
