package models

import (
	"time"
)

type MerchantAccountDetail struct {
	// gorm.Model
	MerchantID    uint
	ApiKey        string    `gorm:"size:255"`
	SecretKey     string    `gorm:"size:255"`
	AccountNumber string    `gorm:"size:255"`
	IfscCode      string    `gorm:"size:244"`
	Merchant      *Merchant `gorm:"foreignKey:MerchantID"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
