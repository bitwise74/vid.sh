// Package model defines database models
package model

type File struct {
	ID      uint   `gorm:"primaryKey;autoIncrement;index" json:"id"`
	UserID  string `json:"-"`
	FileKey string `json:"file_key"` // Avoids file name conflicts
	// ThumbKey     string      `json:"thumb_key"` // TODO: drop this column its not mandatory
	OriginalName string `json:"name"` // Original file name before turning it into a special S3 key
	Private      bool   `json:"private"`
	Format       string `json:"format"`
	// Views        int32       `json:"views"` // TODO: drop
	Size      int64       `json:"size"`
	Tags      StringSlice `json:"tags"`
	State     string      `json:"state"` // Used to inform the frontend/backend if the file is being processed/uploaded
	Version   int         `gorm:"default:1" json:"version"`
	Duration  float64     `json:"duration"` // All are unix millisecond timestamps
	CreatedAt int64       `gorm:"not null" json:"created_at"`
	ExpiresAt *int64      `json:"expires_at,omitzero"`
}
