package model

import "time"

type BlacklistedToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time
}
