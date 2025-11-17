package model

import (
	"salehelp/domain"
	"time"

	"gorm.io/gorm"
)

type Merchant struct {
	MerchantUID string         `gorm:"primaryKey;column:merchant_uid;type:varchar(100)"`
	CompanyName string         `gorm:"column:company_name;type:varchar(512)"`
	IsEnabled   bool           `gorm:"column:is_enabled"`
	Login       string         `gorm:"column:login;type:varchar(100)"`
	Password    string         `gorm:"column:password;type:varchar(100)"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func MerchantColumns() []string {
	return []string{
		"merchant_uid",
		"company_name",
		"login",
		"password",
		"is_enabled",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (m *Merchant) ToDomain() *domain.Merchant {
	if m == nil {
		return nil
	}
	return &domain.Merchant{
		MerchantUID: m.MerchantUID,
		CompanyName: m.CompanyName,
		IsEnabled:   m.IsEnabled,
		Login:       m.Login,
		Password:    m.Password,
	}
}

func (*Merchant) TableName() string {
	return "merchants"
}
