package models

import (
	"time"

	"gorm.io/gorm"
)

type VerifyEmail struct {
	gorm.Model
	UserID     string    `json:"user_id"`
	User       User      `gorm:"references:Username;constraint:OnDelete:SET NULL" json:"user"`
	Code       string    `gorm:"not null" json:"code"`
	Expiration time.Time `gorm:"not null" json:"expiration"`
	IsUsed     bool      `gorm:"default:false" json:"is_used"`
}
