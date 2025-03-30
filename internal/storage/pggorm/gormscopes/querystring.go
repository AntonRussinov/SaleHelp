package gormscopes

import (
	"context"
	"database/sql"
	"salehelp/internal/storage/params/querystring"

	"gorm.io/gorm"
)

func QueryString(ctx context.Context, queryStringTmpl string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		qs := querystring.GetFromContext(ctx)
		if qs == "" {
			return db
		}

		return db.Where(queryStringTmpl, sql.Named("queryString", qs))
	}
}
