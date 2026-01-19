package subscriptionaddon_test

import (
	"log/slog"
	"testing"
	"time"

	"salehelp/domain"
	subscriptionaddon "salehelp/internal/storage/pggorm/repository/subscription_addon"

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

func setupAddonRepo(t *testing.T) (*subscriptionaddon.Repository, sqlmock.Sqlmock) {
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
	repo := subscriptionaddon.NewRepository(log, gormDB)

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

func (b *DataBuilder) AddonTelegramExtraChats(subscriptionID int, qty int) domain.SubscriptionAddon {
	b.t.Helper()

	return domain.SubscriptionAddon{
		SubscriptionID: subscriptionID,
		Code:           "telegram_extra_chat",
		Quantity:       qty,
		AmountMinor:    int64(qty) * 500 * 100 * 1, // пример: 500 RUB/мес/чат, период 1 месяц
		CurrencyCode:   "RUB",
		Meta:           []byte(`{}`),
		CreatedAt:      b.now,
	}
}
