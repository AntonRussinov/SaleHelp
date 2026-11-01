package cronjob_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCronJobRepositoryGetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	tests := []struct {
		name    string
		found   bool
		sqlErr  error
		wantErr bool
	}{
		{name: "returns job", found: true},
		{name: "returns nil when not found", found: false},
		{name: "returns error on sql error", sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			job := builder.CronJob()
			job.ID = 1

			rows := sqlmock.NewRows(JoinedColumns())
			if tt.found {
				rows.AddRow(builder.CronJobRowValues(job)...)
			}

			mock.ExpectQuery(`SELECT (.+) FROM "cron_jobs" WHERE id = \$1(.+)`).
				WithArgs(job.ID, 1).
				WillReturnRows(rows).
				WillReturnError(tt.sqlErr)

			got, err := repo.GetByID(ctx, job.ID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if !tt.found {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, job.ID, got.ID)
			require.Equal(t, job.Name, got.Name)
		})
	}
}

func TestCronJobRepositoryGetLastByName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	name := "snapshot_sync"

	job := builder.CronJob()
	job.Name = name

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.CronJobRowValues(job)...)

	mock.ExpectQuery(`SELECT (.+) FROM "cron_jobs" WHERE name = \$1 ORDER BY started_at DESC LIMIT \$2`).
		WithArgs(name, 1).
		WillReturnRows(rows)

	got, err := repo.GetLastByName(ctx, name)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, name, got.Name)
}

func TestCronJobRepositoryListByName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	builder := NewDataBuilder(t)

	repo, mock := setupRepo(t)

	name := "snapshot_sync"

	j1 := builder.CronJob()
	j1.Name = name
	j2 := builder.CronJob()
	j2.Name = name

	rows := sqlmock.NewRows(JoinedColumns()).
		AddRow(builder.CronJobRowValues(j1)...).
		AddRow(builder.CronJobRowValues(j2)...)

	mock.ExpectQuery(`SELECT (.+) FROM "cron_jobs" WHERE name = \$1 ORDER BY started_at DESC`).
		WithArgs(name).
		WillReturnRows(rows)

	list, err := repo.ListByName(ctx, name, 0, 0)

	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, name, list[0].Name)
	require.Equal(t, name, list[1].Name)
}
