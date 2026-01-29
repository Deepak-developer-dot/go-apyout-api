package models

import (
	"time"
)

type CustomerDetails struct {
	ID uint `gorm:"primaryKey"`
	// gorm.Model
	TransactionID string
	Name          string `gorm:"size:255"`
	Mobile        string `gorm:"size:255"`
	Email         string `gorm:"size:255"`
	AccountNumber string `gorm:"size:255"`
	IfscCode      string `gorm:"size:244"`
	// Transaction   *WalletTransaction `gorm:"foreignKey:TransactionID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (CustomerDetails) TableName() string {
	return "customer_details"
}
