package subscription_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRepositoryCreate(t *testing.T) {
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

			repo, mock := setupSubscriptionRepo(t)
			sub := builder.SubscriptionV0() // сделай билдер под новую структуру

			args := saveSubscriptionV0SQLArgs(sub)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "subscriptions"(.+)`).
				WithArgs(args...).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			subParam := sub
			err := repo.Create(ctx, &subParam)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveSubscriptionV0SQLArgs(s domain.Subscription) []driver.Value {
	// id bigserial: gorm его не вставляет
	// Важно: порядок должен соответствовать тому, как gorm строит insert по полям struct.
	var paidAt any = s.PaidAt
	var comment any = s.Comment
	// source not null
	if s.Source == "" {
		s.Source = "manual"
	}
	return []driver.Value{
		s.ShopID,
		s.Plan,
		s.Status,
		s.ValidFrom,
		s.ValidTo,
		s.PeriodMonths,
		s.BaseAmountMinor,
		s.CurrencyCode,
		paidAt,
		s.Source,
		comment,
		s.CreatedAt,
		s.UpdatedAt,
	}
}

// пример, если нужно
func ptrTime(t time.Time) *time.Time { return &t }
