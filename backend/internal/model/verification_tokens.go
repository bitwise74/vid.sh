package model

import "time"

type VerificationToken struct {
	ID        int    `gorm:"primaryKey;autoincrement"`
	UserID    string `gorm:"index"`
	Token     string `gorm:"uniqueIndex"`
	Purpose   string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
	CleanupAt *time.Time
	Used      bool
}
