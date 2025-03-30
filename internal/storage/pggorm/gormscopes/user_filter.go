package gormscopes

import (
	"context"
	"salehelp/domain"
	"salehelp/pkg/auth"

	"gorm.io/gorm"
)

func UserFilter(ctx context.Context, condition string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		user := auth.GetUserFromContext(ctx)
		if user == nil {
			_ = db.AddError(domain.NewForbiddenErr("unauthorized"))
			return db
		}

		return db.Where(condition, user.ID)
	}
}
