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

// setupSubscriptionRepo поднимает gorm.DB поверх sqlmock
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

func (b *DataBuilder) SubscriptionV0() domain.Subscription {
	b.t.Helper()

	validFrom := b.now.Add(-time.Hour)
	validTo := b.now.Add(30 * 24 * time.Hour)

	paidAt := b.now.Add(-30 * time.Minute)
	comment := "manual grant"

	return domain.Subscription{
		// ID не задаём — BIGSERIAL
		ShopID:          123,
		Plan:            "basic",
		Status:          "active",
		ValidFrom:       validFrom,
		ValidTo:         validTo,
		PeriodMonths:    1,
		BaseAmountMinor: 1500 * 100, // 1500.00 в minor units
		CurrencyCode:    "RUB",
		PaidAt:          &paidAt,
		Source:          "manual",
		Comment:         &comment,
		CreatedAt:       b.now,
		UpdatedAt:       b.now,
	}
}
