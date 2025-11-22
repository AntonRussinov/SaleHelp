package cronjob_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCronJobRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		sqlErr  error
		wantErr bool
	}{
		{name: "executes insert query"},
		{name: "returns error on sql insert error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)
			job := builder.CronJob()

			args := saveCronJobSQLArgs(job)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "cron_jobs"(.+)`).
				WithArgs(args...).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			jobParam := job

			err := repo.Create(ctx, &jobParam)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveCronJobSQLArgs(c domain.CronJob) []driver.Value {
	// id — BIGSERIAL, не передаём
	return []driver.Value{
		c.Name,
		c.Status,
		c.ErrorMessage,
		c.StartedAt,
		c.FinishedAt,
	}
}
