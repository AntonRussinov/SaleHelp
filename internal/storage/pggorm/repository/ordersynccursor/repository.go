package ordersynccursor

import (
	"context"
	"log/slog"
	"time"

	"salehelp/internal/storage/pggorm"
	"salehelp/internal/storage/pggorm/model"
	"salehelp/pkg/logger/sl"

	"gorm.io/gorm"
)

type Repository struct {
	pggorm.Repository
}

func NewRepository(log *slog.Logger, db *gorm.DB) *Repository {
	r := &Repository{}
	r.Log = log.With(sl.Module("OrderSyncCursorRepository"))
	r.DB = db
	return r
}

// Get возвращает last_change_date для пагинации Orders API. Zero time, если курсор не задан.
func (r *Repository) Get(ctx context.Context, shopID int) (time.Time, error) {
	var m model.OrderSyncCursor
	err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	return m.LastChangeDate, nil
}

// Set сохраняет курсор для магазина (upsert по shop_id).
func (r *Repository) Set(ctx context.Context, shopID int, lastChangeDate time.Time) error {
	var m model.OrderSyncCursor
	err := r.DB.WithContext(ctx).Where("shop_id = ?", shopID).First(&m).Error
	if err == gorm.ErrRecordNotFound {
		return r.DB.WithContext(ctx).Create(&model.OrderSyncCursor{ShopID: shopID, LastChangeDate: lastChangeDate}).Error
	}
	if err != nil {
		return err
	}
	m.LastChangeDate = lastChangeDate
	return r.DB.WithContext(ctx).Save(&m).Error
}
