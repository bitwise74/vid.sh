package model

import "time"

type Migration struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}
