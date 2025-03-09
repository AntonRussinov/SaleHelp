package order

import (
	"context"

	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

type (
	key string

	Direction string

	Param struct {
		Field string
		Order Direction
	}
)

const (
	Asc  Direction = "asc"
	Desc Direction = "desc"

	paramsKey key = "orderParams"
)

func GetFromContext(ctx context.Context) []Param {
	order, ok := ctx.Value(paramsKey).([]Param)
	if !ok {
		return nil
	}

	return order
}

func SetToContext(ctx context.Context, params []Param) context.Context {
	return context.WithValue(ctx, paramsKey, params)
}

func GetPaginatorRules(params []Param, fields map[string]string, modelNames map[string]string) []paginator.Rule {
	if len(params) == 0 {
		return nil
	}

	rules := make([]paginator.Rule, 0, len(params))

	for _, param := range params {
		name, ok := fields[param.Field]
		if !ok {
			continue
		}

		if name == "" {
			name = param.Field
		}

		var direction paginator.Order
		if param.Order == Desc {
			direction = paginator.DESC
		} else {
			direction = paginator.ASC
		}

		rule := paginator.Rule{
			Key:     modelNames[param.Field],
			SQLRepr: name,
			Order:   direction,
		}

		rules = append(rules, rule)
	}

	return rules
}
