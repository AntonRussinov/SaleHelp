package gormscopes_test

import (
	"context"
	"salehelp/internal/storage/params/order"
	"salehelp/internal/storage/pggorm/gormscopes"
	"salehelp/test/dbmock"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	t.Parallel()

	orderMap := map[string]string{
		"name": "name",
		"type": "templates.type",
		"date": "",
	}

	tests := []struct {
		name            string
		params          []order.Param
		wantEqual       bool
		predefinedOrder string
	}{
		{name: "no params", params: nil, wantEqual: true},
		{name: "apply params matching map entries", params: []order.Param{{
			Field: "name", Order: order.Asc,
		}, {Field: "type", Order: order.Desc}}, wantEqual: false},
		{name: "not apply missing in map params", params: []order.Param{{
			Field: "missing", Order: order.Asc,
		}}, wantEqual: true},
		{name: "apply predefined order", params: nil, wantEqual: false, predefinedOrder: "name"},
		{name: "apply params if map value missing", params: []order.Param{{
			Field: "date", Order: order.Asc,
		}}, wantEqual: false, predefinedOrder: "name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, _ := dbmock.MockGORM(t)

			ctx := context.Background()
			if tt.params != nil {
				ctx = order.SetToContext(ctx, tt.params)
			}

			scopeFunc := gormscopes.Order(ctx, orderMap, tt.predefinedOrder)

			res := scopeFunc(db)
			require.NoError(t, res.Error)

			if tt.wantEqual {
				require.Equal(t, res, db)
			} else {
				require.NotEqual(t, res, db)
			}
		})
	}
}
