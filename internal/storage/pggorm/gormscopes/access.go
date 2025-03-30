package gormscopes

import (
	"context"
	"salehelp/domain"
	"salehelp/pkg/auth"

	"gorm.io/gorm"
)

func AccessFilterWithTableName(ctx context.Context, table string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		user := auth.GetUserFromContext(ctx)
		if user == nil {
			_ = db.AddError(domain.NewForbiddenErr("unauthorized"))
			return db
		}

		if user.HasRole(domain.UserRoleAdmin) || user.HasRole(domain.UserRoleSuperAdmin) {
			return db
		}

		if table == "" {
			return db.Where("author_id = ?", user.ID)
		}

		return db.Where(table+".author_id = ?", user.ID)
	}
}

func AccessFilter(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return AccessFilterWithTableName(ctx, "")
}
