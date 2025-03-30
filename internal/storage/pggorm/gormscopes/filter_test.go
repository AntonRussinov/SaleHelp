package gormscopes_test

import (
	"context"
	"salehelp/internal/storage/params/filter"
	"salehelp/internal/storage/pggorm/gormscopes"
	"salehelp/test/dbmock"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	filterMap := gormscopes.FilterMap{
		"app_id": {
			FieldName:   "app_id",
			ValueType:   filter.ValueTypeString,
			CompareType: filter.CompareTypeEqual,
		},
		"app_name": {
			FieldName: "app_name",
			ValueType: filter.ValueTypeString,
		},
		"app_status": {
			ValueType: filter.ValueTypeString,
		},
		"int_field": {
			FieldName: "int_field",
			ValueType: filter.ValueTypeInt,
		},
		"array_field": {
			FieldName:   "array_field",
			ValueType:   filter.ValueTypeFloat,
			CompareType: filter.CompareTypeIn,
		},
		"bool_field": {
			FieldName: "bool_field",
			ValueType: filter.ValueTypeBool,
		},
		"uuid_field": {
			FieldName: "uuid_field",
			ValueType: filter.ValueTypeUUID,
		},
		"date_field": {
			FieldName: "date_field",
			ValueType: filter.ValueTypeDate,
		},
		"float_field": {
			FieldName: "float_field",
			ValueType: filter.ValueTypeFloat,
		},
		"no_value_type": {
			FieldName: "no_value_type",
		},
		"tag": {
			FieldName:   "display_tags",
			ValueType:   filter.ValueTypeString,
			CompareType: filter.CompareTypeLike,
		},
	}

	tests := []struct {
		name        string
		paramsInCtx []filter.Param
		wantEqual   bool
	}{
		{name: "no params", paramsInCtx: nil, wantEqual: true},
		{name: "zero params", paramsInCtx: []filter.Param{}, wantEqual: true},
		{name: "param not in map", paramsInCtx: []filter.Param{
			{Field: "missing", Value: "value", CompareType: filter.CompareTypeEqual},
		}, wantEqual: true},
		{name: "missing name in map fallback to name in query", paramsInCtx: []filter.Param{
			{Field: "app_status", Value: "value"},
		}, wantEqual: false},
		{name: "compare type in", paramsInCtx: []filter.Param{
			{Field: "app_id", Values: []string{"1", "2"}, CompareType: filter.CompareTypeIn},
		}, wantEqual: false},
		{name: "fail to covert value", paramsInCtx: []filter.Param{
			{Field: "int_field", Value: "ABC", CompareType: filter.CompareTypeEqual},
		}, wantEqual: true},
		{name: "compare type in", paramsInCtx: []filter.Param{
			{Field: "array_field", Values: []string{"1", "2"}, CompareType: filter.CompareTypeIn},
		}, wantEqual: false},
		{name: "seacrh by tag", paramsInCtx: []filter.Param{
			{Field: "tag", Values: []string{"tag1", "tag2"}, CompareType: filter.CompareTypeEqual, Value: "tag1,tag2"},
		}, wantEqual: false},
		{name: "different value types", paramsInCtx: []filter.Param{
			{Field: "int_field", Value: "1", CompareType: filter.CompareTypeEqual},
			{Field: "bool_field", Value: "true", CompareType: filter.CompareTypeEqual},
			{Field: "uuid_field", Value: "00000000-0000-0000-0000-000000000001", CompareType: filter.CompareTypeEqual},
			{Field: "date_field", Value: "2020-01-01", CompareType: filter.CompareTypeEqual},
			{Field: "float_field", Value: "1.1", CompareType: filter.CompareTypeEqual},
			{Field: "no_value_type", Value: "value", CompareType: filter.CompareTypeEqual},
		}, wantEqual: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, _ := dbmock.MockGORM(t)

			ctx := context.Background()
			if tt.paramsInCtx != nil {
				ctx = filter.SetToContext(ctx, tt.paramsInCtx)
			}

			scopeFunc := gormscopes.Filter(ctx, filterMap)
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
