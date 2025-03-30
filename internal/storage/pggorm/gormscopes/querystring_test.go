package gormscopes_test

import (
	"context"
	"salehelp/internal/storage/params/querystring"
	"salehelp/internal/storage/pggorm/gormscopes"
	"salehelp/test/dbmock"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		param     string
		wantEqual bool
	}{
		{name: "param is applied", param: "app-1", wantEqual: false},
		{name: "missing params not applied", param: "", wantEqual: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, _ := dbmock.MockGORM(t)
			ctx := context.Background()
			ctx = querystring.SetToContext(ctx, tt.param)
			scopeFunc := gormscopes.QueryString(ctx, "LIKE '%%' || :queryString || '%%'")

			res := scopeFunc(db)
			if tt.wantEqual {
				require.Equal(t, res, db)
			} else {
				require.NotEqual(t, res, db)
			}
		})
	}
}
