package model

import (
	"encoding"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Compile-time checks to ensure User implements the encoding interfaces
var (
	_ encoding.BinaryMarshaler   = (*User)(nil)
	_ encoding.BinaryUnmarshaler = (*User)(nil)
)

type User struct {
	ID           string `gorm:"primaryKey;autoIncrement" json:"-"`
	Email        string `gorm:"unique; not null" json:"-"`
	PasswordHash string `gorm:"not null" json:"-"`
	Verified     bool   `gorm:"default:false" json:"-"`

	ExpiresAt *time.Time     `json:"-"`
	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`

	AvatarHash           string `gorm:"default:null" json:"avatarHash"`
	Username             string `gorm:"default:null;unique" json:"username"`
	PublicProfileEnabled bool   `gorm:"default:false" json:"publicProfileEnabled"`
	DefaultPrivateVideos bool   `gorm:"default:true" json:"defaultPrivateVideos"`

	// VerificationTokens []VerificationToken `gorm:"foreignKey:UserID" json:"-"` // TODO: drop this
	Tokens []Token `gorm:"foreignKey:UserID" json:"-"`
	Files  []File  `gorm:"foreignKey:UserID" json:"videos"` // TODO: rename to Videos
	Stats  Stats   `gorm:"foreignKey:UserID" json:"stats"`
	// ResendRequests ResendRequest `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
