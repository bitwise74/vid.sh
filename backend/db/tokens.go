package db

import (
	"bitwise74/video-api/internal/model"
	"time"
)

// BlacklistJWTToken adds a JWT token to the blacklist with an expiry time.
func (d *DB) BlacklistJWTToken(token string, expiry time.Time) error {
	return d.Gorm.Create(&model.BlacklistedToken{
		Token:     token,
		ExpiresAt: expiry,
	}).Error
}

func (d *DB) SaveToken(t *model.Token) error {
	return d.Gorm.Create(t).Error
}
