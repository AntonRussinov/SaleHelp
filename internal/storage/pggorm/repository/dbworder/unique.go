package dbworder

import (
	"github.com/jackc/pgx/v5/pgconn"
)

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if err == nil {
		return false
	}
	if !asPgError(err, &pgErr) {
		return false
	}
	return pgErr.Code == "23505" // unique_violation
}

func asPgError(err error, out **pgconn.PgError) bool {
	for err != nil {
		if e, ok := err.(*pgconn.PgError); ok {
			*out = e
			return true
		}
		switch e := err.(type) {
		case interface{ Unwrap() error }:
			err = e.Unwrap()
		default:
			return false
		}
	}
	return false
}
