package model

import (
	"salehelp/domain"
	"time"

	"gorm.io/gorm"
)

// Image представляет таблицу images.
type Image struct {
	ID        int            `gorm:"primaryKey;column:id"`
	OfferID   int            `gorm:"column:offer_id"`
	ImagePath string         `gorm:"column:image_path"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func ImageColumns() []string {
	return []string{
		"id",
		"offer_id",
		"image_path",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (img *Image) ToDomain() *domain.Image {
	if img == nil {
		return nil
	}
	return &domain.Image{
		ID:        img.ID,
		OfferID:   img.OfferID,
		ImagePath: img.ImagePath,
	}
}

func (*Image) TableName() string {
	return "images"
}
