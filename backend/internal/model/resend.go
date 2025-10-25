package model

import "time"

// TODO: drop table
type ResendRequest struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	UserID     string
	LastResend time.Time
	Cooldown   time.Time
	Blocked    bool // If the user sends too many resend requests they're blocked for the day
}
