package pggorm

import (
	"errors"
	"fmt"
	"log/slog"
	"salehelp/internal/storage"
	"salehelp/pkg/logger/sl"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository struct {
	Log *slog.Logger
	DB  *gorm.DB
}

func (repo *Repository) HandleError(fn string, err error) error {
	repo.Log.Error("error in repository: ", sl.Method(fn), sl.Err(err))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%s: %w", fn, storage.ErrNotFound)
	}

	var pgerr *pgconn.PgError

	if errors.As(err, &pgerr) {
		return parsePgConnError(fn, pgerr)
	}

	return wrapError(fn, err)
}

func GroupQuery(groupBy string) string {
	groupByMap := map[string]string{
		"days":   "DATE(created_at)",
		"weeks":  "DATE_TRUNC('week', created_at)",
		"months": "DATE_TRUNC('month', created_at)",
		"years":  "DATE_TRUNC('year', created_at)",
		"60days": "DATE_TRUNC('day', created_at) - ((EXTRACT(DAY FROM created_at)::integer - 1) % 60) * INTERVAL '1 day'",
		"90days": "DATE_TRUNC('day', created_at) - ((EXTRACT(DAY FROM created_at)::integer - 1) % 90) * INTERVAL '1 day'",
	}

	return groupByMap[strings.ToLower(groupBy)]
}
