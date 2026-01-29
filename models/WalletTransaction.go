package models

import (
	"time"
)

type WalletTransaction struct {
	ID uint `gorm:"primaryKey"`
	// gorm.Model
	MerchantID         uint    `gorm:"size:255"`
	CreatedBy          uint    `gorm:"size:255"`
	Amount             float64 `gorm:"size:255"`
	PayoutProvider     uint    `gorm:"size:255"`
	ProviderCommission float64 `gorm:"size:255"`
	CompanyCommission  float64 `gorm:"size:255"`
	WithdrawCommission float64 `gorm:"size:255"`
	GstCharge          float64 `gorm:"size:255"`
	ServiceCharge      float64 `gorm:"size:255"`
	PaymentType        uint    `gorm:"size:255"`
	PaymentTypeString  string  `gorm:"size:255"`

	TransactionID         string    `gorm:"size:255"`
	ClientTransactionID   string    `gorm:"size:255"`
	ProviderTransactionID string    `gorm:"size:255"`
	Utr                   string    `gorm:"size:255"`
	PType                 uint      `gorm:"size:20"`
	CurWalletBal          float64   `gorm:"size:255"`
	PaymentMode           string    `gorm:"size:255"`
	PlatformMode          string    `gorm:"size:255"`
	PaymentDetails        string    `gorm:"size:255"`
	CustomerDetails       string    `gorm:"size:255"`
	ContactMobileNumber   string    `gorm:"size:255"`
	Status                uint      `gorm:"size:100"`
	StatusNote            string    `gorm:"size:100"`
	IsWebhookSent         uint      `gorm:"size:100"`
	IsSelfWithdrawal      uint      `gorm:"size:100"`
	LocationDetails       string    `gorm:"size:100"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	CustomerDetail CustomerDetails `gorm:"foreignKey:TransactionID;references:TransactionID"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}
