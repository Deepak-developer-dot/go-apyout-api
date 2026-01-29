package models

import "time"

type Wallet struct {
	ID         uint      `gorm:"primaryKey"`
	Amount     float64   `gorm:"amount"`
	MerchantID uint      `gorm:"merchant_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Wallet) TableName() string {
	return "wallet"
}
