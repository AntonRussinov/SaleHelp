package gormscopes

import (
	"context"
	"salehelp/internal/storage/params/order"
	"strings"

	"gorm.io/gorm"
)

func Order(ctx context.Context, orderMap map[string]string, def string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		statemet := db

		orderFields := order.GetFromContext(ctx)
		if len(orderFields) > 0 {
			orderString := getOrderString(orderFields, orderMap)
			if orderString != "" {
				statemet = db.Order(orderString)
			}
		}

		if def != "" {
			statemet = statemet.Order(def)
		}

		return statemet
	}
}

func getOrderString(params []order.Param, orderMap map[string]string) string {
	orderStrings := make([]string, 0)

	for _, param := range params {
		name, ok := orderMap[param.Field]
		if !ok {
			continue
		}

		if name == "" {
			name = param.Field
		}

		direction := string(param.Order)

		orderString := name + " " + direction
		if param.Order == order.Desc {
			orderString += " NULLS LAST"
		} else {
			orderString += " NULLS FIRST"
		}

		orderStrings = append(orderStrings, orderString)
	}

	return strings.Join(orderStrings, ", ")
}
