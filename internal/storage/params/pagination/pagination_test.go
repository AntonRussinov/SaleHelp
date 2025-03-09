package pagination_test

import (
	"context"
	"salehelp/internal/storage/params/pagination"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasNext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		params     *pagination.Params
		totalCount int64
		expected   bool
	}{
		{name: "no params", params: nil, totalCount: 1200, expected: false},
		{name: "zero limit", params: &pagination.Params{Limit: 0, Offset: 12}, totalCount: 1200, expected: false},
		{name: "all items requested", params: &pagination.Params{Limit: 10, Offset: 0}, totalCount: 10, expected: false},
		{name: "has items left", params: &pagination.Params{Limit: 10, Offset: 0}, totalCount: 11, expected: true},
		{name: "requested more items", params: &pagination.Params{Limit: 200, Offset: 200}, totalCount: 345, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := pagination.SetToContext(context.Background(), tt.params)
			res := pagination.HasNext(ctx, tt.totalCount)
			assert.Equal(t, tt.expected, res)
		})
	}
}
