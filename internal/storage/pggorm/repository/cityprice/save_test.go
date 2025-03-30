package cityprice_test

import (
	"context"
	"database/sql/driver"
	"salehelp/domain"
	"salehelp/test/dbmock"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCityPriceRepositorySave(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	builder := NewDataBuilder(t)

	tests := []struct {
		name      string
		cityPrice domain.CityPrice
		sqlErr    error
		wantErr   bool
	}{
		{name: "executes insert query", cityPrice: builder.CityPrice()},
		{name: "returns error on sql insert error", cityPrice: builder.CityPrice(), sqlErr: assert.AnError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := setupRepo(t)

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO "cityPrices"(.+)`).WithArgs(saveCityPriceSQLArgs(tt.cityPrice)...).
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.sqlErr)

			if tt.sqlErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			cityPriceParam := tt.cityPrice

			err := repo.Save(ctx, &cityPriceParam)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func saveCityPriceSQLArgs(cityPrice domain.CityPrice) []driver.Value {
	values := []driver.Value{
		cityPrice.ID,
		cityPrice.MerchantUID,
		cityPrice.OfferID,
		cityPrice.CityID,
		cityPrice.Price,
		cityPrice.OldPrice,
		dbmock.MatchTime(),
		dbmock.MatchTime(),
		nil,
	}

	return values
}
