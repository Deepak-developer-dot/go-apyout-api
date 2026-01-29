package models

import (
	"time"
)

type Merchant struct {
	// gorm.Model
	ID                     uint   `gorm:"primaryKey;autoIncrement"`
	OwnerName              string `gorm:"size:255"`
	Email                  string `gorm:"uniqueIndex;size:255"`
	Password               string `gorm:"size:255"`
	ProviderType           uint   `gorm:"size:50"`
	IsAdmin                uint   `gorm:"size:50"`
	Mobile                 uint   `gorm:"size:20"`
	AppName                string `gorm:"size:255"`
	AppUrl                 string `gorm:"size:255"`
	WebUrl                 string `gorm:"size:255"`
	WebhookUrl             string `gorm:"size:255"`
	IsWebhookActive        uint   `gorm:"size:50"`
	IpAddress              string `gorm:"size:100"`
	WithdrawalCommission   float64
	CommissionChargeType   uint `gorm:"size:100"`
	Status                 uint `gorm:"size:20"`
	IsDeleted              uint `gorm:"size:20"`
	IsDependent            uint `gorm:"size:20"`
	AutomaticCreditEnabled uint `gorm:"size:20"`
	DirectPayoutEnabled    uint `gorm:"size:20"`
	AgentID                uint `gorm:"size:20"`
	// PartnerID              uint      `gorm:"size:20"`
	// TypeSuperAdmin      uint      `gorm:"size:20"`
	PayoutProvider      uint `gorm:"size:20"`
	PerTransactionLimit uint `gorm:"size:20"`
	// SortBy              uint      `gorm:"size:20"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	MerchantAccountDetail MerchantAccountDetail `gorm:"foreignKey:MerchantID"`
}

func (Merchant) TableName() string {
	return "admins"
}
