package gormscopes

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/pkg/auth"

	"gorm.io/gorm"
)

func ApplicationFiltersWithTableName(ctx context.Context, table string) func(db *gorm.DB) *gorm.DB {
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
			return db.Where(`
			(author_id = ? 
			OR access_mode = 'global' 
			OR EXISTS(
				SELECT usr.user_id 
					FROM application_users usr 
				WHERE usr.application_id = applications.id
					AND usr.user_id = ?
			)
			OR EXISTS(
				SELECT fnusers.user_id
					FROM function_group_users fnusers
				WHERE fnusers.user_id = ?
					AND fnusers.function_group_id IN (
						SELECT fns.function_group_id 
							FROM application_frunction_groups fns 
						WHERE fns.application_id = applications.id
				)
			)
			)`, user.ID, user.ID, user.ID)
		}

		condition := fmt.Sprintf(`(
		%s.author_id = ? 
		OR %s.access_mode = 'global' 
		OR EXISTS(
			SELECT usr.user_id 
				FROM application_users usr 
			WHERE usr.application_id = %s.id
				AND usr.user_id = ?
		)
		OR EXISTS(
			SELECT fnusers.user_id
					FROM function_group_users fnusers
				WHERE fnusers.user_id = ?
					AND fnusers.function_group_id IN (
						SELECT fns.function_group_id 
							FROM application_function_groups fns 
						WHERE fns.application_id = %s.id
				)
		)
		)`, table, table, table, table)

		return db.Where(condition, user.ID, user.ID, user.ID)
	}
}

func ApplicationFilters(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return ApplicationFiltersWithTableName(ctx, "")
}
