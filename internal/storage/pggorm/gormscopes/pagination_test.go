package gormscopes_test

import (
	"context"
	"salehelp/internal/storage/params/pagination"
	"salehelp/internal/storage/pggorm/gormscopes"
	"salehelp/test/dbmock"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPagination(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mustPaginate bool
		params       *pagination.Params
		wantErr      bool
	}{
		{name: "must paginate no params", mustPaginate: true, wantErr: true},
		{name: "no params", mustPaginate: false, wantErr: false},
		{name: "apply params", mustPaginate: false, params: &pagination.Params{Limit: 10, Offset: 20}, wantErr: false},
		{name: "must apply params", mustPaginate: true, params: &pagination.Params{Limit: 10, Offset: 20}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, _ := dbmock.MockGORM(t)

			ctx := context.Background()
			if tt.params != nil {
				ctx = pagination.SetToContext(ctx, tt.params)
			}

			var scopeFunc func(db *gorm.DB) *gorm.DB
			if tt.mustPaginate {
				scopeFunc = gormscopes.MustPagination(ctx)
			} else {
				scopeFunc = gormscopes.Pagination(ctx)
			}

			res := scopeFunc(db)
			if tt.wantErr {
				require.Error(t, res.Error)
			} else {
				require.NoError(t, res.Error)
			}
		})
	}
}
