package order_test

import (
	"salehelp/internal/storage/params/order"
	"testing"

	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"github.com/stretchr/testify/assert"
)

func TestGetPaginatorRules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		params     []order.Param
		fields     map[string]string
		modelNames map[string]string
		want       []paginator.Rule
	}{
		{
			name:   "no params",
			params: []order.Param{},
			fields: map[string]string{
				"name": "applications.name",
				"id":   "applications.id",
			},
			modelNames: map[string]string{
				"name": "Name",
				"id":   "ID",
			},
			want: nil,
		},
		{
			name: "params with matching fields and modelNames",
			params: []order.Param{
				{Field: "name", Order: order.Asc},
				{Field: "id", Order: order.Desc},
			},
			fields: map[string]string{
				"name": "applications.name",
				"id":   "applications.id",
			},
			modelNames: map[string]string{
				"name": "Name",
				"id":   "ID",
			},
			want: []paginator.Rule{
				{Key: "Name", SQLRepr: "applications.name", Order: paginator.ASC},
				{Key: "ID", SQLRepr: "applications.id", Order: paginator.DESC},
			},
		},
		{
			name: "params with missing fields in fields map",
			params: []order.Param{
				{Field: "missing_field", Order: order.Asc},
				{Field: "id", Order: order.Desc},
			},
			fields: map[string]string{
				"id": "applications.id",
			},
			modelNames: map[string]string{
				"id": "ID",
			},
			want: []paginator.Rule{
				{Key: "ID", SQLRepr: "applications.id", Order: paginator.DESC},
			},
		},
		{
			name: "params with missing fields in modelNames map",
			params: []order.Param{
				{Field: "name", Order: order.Asc},
				{Field: "missing_field", Order: order.Desc},
			},
			fields: map[string]string{
				"name": "applications.name",
			},
			modelNames: map[string]string{
				"name": "Name",
			},
			want: []paginator.Rule{
				{Key: "Name", SQLRepr: "applications.name", Order: paginator.ASC},
			},
		},
		{
			name: "params with empty fields in fields map",
			params: []order.Param{
				{Field: "empty_field", Order: order.Asc},
				{Field: "id", Order: order.Desc},
			},
			fields: map[string]string{
				"empty_field": "",
				"id":          "applications.id",
			},
			modelNames: map[string]string{
				"empty_field": "EmptyField",
				"id":          "ID",
			},
			want: []paginator.Rule{
				{Key: "EmptyField", SQLRepr: "empty_field", Order: paginator.ASC},
				{Key: "ID", SQLRepr: "applications.id", Order: paginator.DESC},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := order.GetPaginatorRules(tt.params, tt.fields, tt.modelNames)
			assert.Equal(t, tt.want, got)
		})
	}
}
