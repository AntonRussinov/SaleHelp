package gormscopes

import (
	"context"
	"errors"
	"salehelp/internal/storage/params/pagination"

	"gorm.io/gorm"
)

var ErrNoPaginationParams = errors.New("no pagination params provided")

func applyParams(db *gorm.DB, params *pagination.Params) *gorm.DB {
	return db.Limit(int(params.Limit)).Offset(int(params.Offset))
}

func MustPagination(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		params := pagination.GetFromContext(ctx)
		if params == nil {
			_ = db.AddError(ErrNoPaginationParams)
			return db
		}

		return applyParams(db, params)
	}
}

func Pagination(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		params := pagination.GetFromContext(ctx)
		if params == nil {
			return db
		}

		return applyParams(db, params)
	}
}
