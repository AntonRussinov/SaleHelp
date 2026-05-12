package subscription_test

import (
	"log/slog"
	"testing"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/repository/subscription"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DataBuilder struct {
	t   *testing.T
	now time.Time
}

func setupSubscriptionRepo(t *testing.T) (*subscription.Repository, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{
			Conn:                 sqlDB,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	require.NoError(t, err)

	log := slog.Default()
	repo := subscription.NewRepository(log, gormDB)

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	return repo, mock
}

func NewDataBuilder(t *testing.T) *DataBuilder {
	t.Helper()
	return &DataBuilder{
		t:   t,
		now: time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC),
	}
}

func (b *DataBuilder) Subscription() domain.Subscription {
	b.t.Helper()
	nextMonth := time.Date(b.now.Year(), b.now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	return domain.Subscription{
		UserID:       42,
		Tier:         domain.TierBusiness,
		Status:       domain.SubStatusActive,
		PeriodMonths: 1,
		StartedAt:    b.now.Add(-time.Hour),
		ExpiresAt:    b.now.Add(30 * 24 * time.Hour),
		AIQuota:      200,
		AIUsed:       0,
		AIResetsAt:   nextMonth,
		MaxShops:     2,
		ExtraShops:   0,
		PriceMinor:   149000,
		CreatedAt:    b.now,
		UpdatedAt:    b.now,
	}
}
