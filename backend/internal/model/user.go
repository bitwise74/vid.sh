package model

import "time"

type User struct {
	ID           string `gorm:"primaryKey;autoIncrement"`
	Email        string `gorm:"unique; not null"`
	PasswordHash string `gorm:"not null"`
	Verified     bool   `gorm:"default:false"`
	ExpiresAt    *time.Time

	VerificationTokens []VerificationToken `gorm:"foreignKey:UserID"`
	Files              []File              `gorm:"foreignKey:UserID"`
	Stats              Stats               `gorm:"foreignKey:UserID"`
	ResendRequests     ResendRequest       `gorm:"foreignKey:UserID"`
}
