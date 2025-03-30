package gormscopes_test

import (
	"context"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/gormscopes"
	"salehelp/pkg/auth"
	"salehelp/test/dbmock"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAccess(t *testing.T) {
	t.Parallel()

	admin := domain.User{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Roles: []domain.UserRole{domain.UserRoleAdmin},
	}
	user := domain.User{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Roles: []domain.UserRole{domain.UserRoleUser},
	}

	tests := []struct {
		name      string
		user      *domain.User
		wantErr   bool
		wantEqual bool
	}{
		{name: "error on no user", wantErr: true, wantEqual: false},
		{name: "work for admin with no change", user: &admin, wantErr: false, wantEqual: true},
		{name: "chage for user", user: &user, wantErr: false, wantEqual: false},
	}

	funccalls := []struct {
		name            string
		fn              func(ctx context.Context) func(db *gorm.DB) *gorm.DB
		fnWithTableName func(ctx context.Context, table string) func(db *gorm.DB) *gorm.DB
		table           string
	}{
		{name: "Access", fn: gormscopes.AccessFilter, table: ""},
		{name: "AccessWithTableName", fn: gormscopes.ApplicationFilters, table: ""},
		{name: "AppFilters", fnWithTableName: gormscopes.AccessFilterWithTableName, table: "applications"},
		{name: "AppFiltersWithTAbleName", fnWithTableName: gormscopes.ApplicationFiltersWithTableName, table: "applications"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, fncall := range funccalls {
				t.Run(fncall.name, func(t *testing.T) {
					db, _ := dbmock.MockGORM(t)
					ctx := context.Background()

					if tt.user != nil {
						ctx = auth.SetUserToContext(ctx, tt.user)
					}

					var scopeFunc func(db *gorm.DB) *gorm.DB
					if fncall.table == "" {
						scopeFunc = fncall.fn(ctx)
					} else {
						scopeFunc = fncall.fnWithTableName(ctx, fncall.table)
					}

					res := scopeFunc(db)

					if tt.wantErr {
						require.Error(t, res.Error)
						return
					}

					if tt.wantEqual {
						require.Equal(t, res, db)
					} else {
						require.NotEqual(t, res, db)
					}
				})
			}
		})
	}
}

func TestAccessFilterWithTableName(t *testing.T) {
	t.Parallel()

	user := domain.User{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Roles: []domain.UserRole{domain.UserRoleUser},
	}

	db, _ := dbmock.MockGORM(t)
	ctx := context.Background()

	ctx = auth.SetUserToContext(ctx, &user)
	scopeFunc := gormscopes.AccessFilterWithTableName(ctx, "applications")
	res := scopeFunc(db)
	require.NoError(t, res.Error)
	require.NotEqual(t, res, db)
}
