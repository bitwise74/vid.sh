package model

import "time"

type Token struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	UserID    string
	Token     string
	Purpose   string
	ExpiresAt time.Time
	CreatedAt time.Time
	CleanupAt *time.Time
	Used      bool
}
