package service

import (
	"bitwise74/video-api/internal/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func StaleTokenCleanup(t time.Duration, d *gorm.DB) {
	zap.L().Debug("Stale token cleanup attached", zap.Duration("tick_every", t))

	ticker := time.NewTicker(t)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			result := d.Model(model.Token{}).
				Where("expires_at < ?", time.Now()).
				Delete(&model.Token{})
			if result.Error != nil {
				zap.L().Error("Failed to delete stale tokens", zap.Error(result.Error))
				continue
			}

			if result.RowsAffected > 0 {
				zap.L().Debug("Deleted stale tokens", zap.Int64("count", result.RowsAffected))
			}
		}
	}()
}

func UnverifiedUserCleanup(t time.Duration, d *gorm.DB) {
	zap.L().Debug("Unverified user cleanup attached", zap.Duration("tick_every", t))

	ticker := time.NewTicker(t)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			result := d.Model(model.User{}).
				Where("expires_at < ? AND verified = ?", time.Now(), false).
				Delete(&model.User{})
			if result.Error != nil {
				zap.L().Error("Failed to delete unverified users", zap.Error(result.Error))
				continue
			}

			if result.RowsAffected > 0 {
				zap.L().Debug("Soft deleted unverified users", zap.Int64("count", result.RowsAffected))
			}
		}
	}()
}
