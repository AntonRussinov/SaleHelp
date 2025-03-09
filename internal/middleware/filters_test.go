package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"salehelp/internal/middleware"
	"salehelp/internal/storage/params/filter"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type FiltersTestCase struct {
	name   string
	query  string
	expect []filter.Param
}

func filtersTestCases() []FiltersTestCase {
	return []FiltersTestCase{
		{name: "empty", query: "", expect: nil},
		{
			name: "one field default compare type", query: "filter[field1]=value1",
			expect: []filter.Param{
				{Field: "field1", CompareType: filter.CompareTypeEqual, Value: "value1", Values: []string{"value1"}},
			},
		},
		{
			name: "one field equal compare type", query: "filter[field1][eq]=value1",
			expect: []filter.Param{
				{Field: "field1", CompareType: filter.CompareTypeEqual, Value: "value1", Values: []string{"value1"}},
			},
		},
		{
			name: "one field not equal compare type", query: "filter[field1][ne]=value1",
			expect: []filter.Param{
				{Field: "field1", CompareType: filter.CompareTypeNotEqual, Value: "value1", Values: []string{"value1"}},
			},
		},
		{
			name:  "compare operator types",
			query: "filter[field1][gt]=value1&filter[field2][gte]=value2&filter[field3][lt]=value3&filter[field4][lte]=value4",
			expect: []filter.Param{
				{Field: "field1", CompareType: filter.CompareTypeGreaterThan, Value: "value1", Values: []string{"value1"}},
				{Field: "field2", CompareType: filter.CompareTypeGreaterOrEqual, Value: "value2", Values: []string{"value2"}},
				{Field: "field3", CompareType: filter.CompareTypeLessThan, Value: "value3", Values: []string{"value3"}},
				{Field: "field4", CompareType: filter.CompareTypeLessOrEqual, Value: "value4", Values: []string{"value4"}},
			},
		},
		{
			name:  "string like compare types",
			query: "filter[field1][like]=value1&filter[field2][nlike]=value2&filter[field3][starts]=value3&filter[field4][ends]=value4",
			expect: []filter.Param{
				{Field: "field1", CompareType: filter.CompareTypeLike, Value: "value1", Values: []string{"value1"}},
				{Field: "field2", CompareType: filter.CompareTypeNotLike, Value: "value2", Values: []string{"value2"}},
				{Field: "field3", CompareType: filter.CompareTypeStartsWith, Value: "value3", Values: []string{"value3"}},
				{Field: "field4", CompareType: filter.CompareTypeEndsWith, Value: "value4", Values: []string{"value4"}},
			},
		},
		{
			name:  "array compare types",
			query: "filter[field1][in]=value1&filter[field2][nin]=value2",
			expect: []filter.Param{
				{Field: "field1", CompareType: filter.CompareTypeIn, Value: "value1", Values: []string{"value1"}},
				{Field: "field2", CompareType: filter.CompareTypeNotIn, Value: "value2", Values: []string{"value2"}},
			},
		},
		{
			name: "invalid compare type", query: "filter[field1][invalid]=value1",
			expect: []filter.Param{
				{Field: "field1", CompareType: filter.CompareTypeEqual, Value: "value1", Values: []string{"value1"}},
			},
		},
		{name: "invalid param prefix", query: "filters[field1][eq]=value1", expect: nil},
		{
			name: "invalid or missing param compare type part", query: "filter[field1][iasdfasdf=value1&filter[field2]=value2",
			expect: []filter.Param{
				{Field: "field1", CompareType: filter.CompareTypeEqual, Value: "value1", Values: []string{"value1"}},
				{Field: "field2", CompareType: filter.CompareTypeEqual, Value: "value2", Values: []string{"value2"}},
			},
		},
		{name: "incorrect filtering field tag", query: "filter[field1[=value", expect: nil},
	}
}

func TestFilters(t *testing.T) {
	t.Parallel()

	tests := filtersTestCases()

	e := echo.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.URL.RawQuery = tt.query
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mw := middleware.Filters()

			handler := mw(func(c echo.Context) error {
				filters := filter.GetFromContext(c.Request().Context())
				if tt.expect == nil {
					require.Nil(t, filters)
				}

				require.ElementsMatch(t, tt.expect, filters)

				return nil
			})

			_ = handler(c)
		})
	}
}
