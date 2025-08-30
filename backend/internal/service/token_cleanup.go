package service

import (
	"bitwise74/video-api/internal/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TokenCleanup defines a function used to periodically cleanup old
// verification tokens that aren't needed anymore
func TokenCleanup(t time.Duration, db *gorm.DB) {
	ticker := time.NewTicker(t)

	zap.L().Debug("Token cleanup attached", zap.Duration("tick_every", t))

	go func() {
		for range ticker.C {
			var toCleanIds []string

			err := db.
				Model(model.VerificationToken{}).
				Where("expires_at < ?", time.Now()).
				Select("id").
				Find(&toCleanIds).
				Error
			if err != nil {
				zap.L().Error("Failed to query db for tokens to clean", zap.Error(err))
				return
			}

			if len(toCleanIds) > 0 {
				zap.L().Debug("Cleaning up expired tokens")

				err = db.
					Where("id LIKE ?", toCleanIds).
					Delete(model.VerificationToken{}).
					Error
				if err != nil {
					zap.L().Error("Failed to cleanup database", zap.Error(err))
				}
			}
		}
	}()
}
