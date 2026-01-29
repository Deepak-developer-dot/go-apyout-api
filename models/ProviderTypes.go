package models

import (
	"time"
)

type ProviderTypes struct {
	// gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProviderTypes) TableName() string {
	return "provider_type"
}
