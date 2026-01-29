package models

import (
	"time"
)

type PayoutProviders struct {
	// gorm.Model
	ID            uint `gorm:"primaryKey"`
	ProviderType  uint
	Name          string         `gorm:"size:255"`
	MerchantName  string         `gorm:"size:255"`
	AppKey        string         `gorm:"size:255"`
	SecretKey     string         `gorm:"size:255"`
	Email         string         `gorm:"size:255"`
	AccountNumber string         `gorm:"size:255"`
	Status        uint           `gorm:"size:244"`
	Balance       float64        `gorm:"size:244"`
	ProviderTypes *ProviderTypes `gorm:"foreignKey:ProviderType"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (PayoutProviders) TableName() string {
	return "payout_providers"
}
