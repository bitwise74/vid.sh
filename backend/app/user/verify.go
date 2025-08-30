package user

import (
	"bitwise74/video-api/internal"
	"bitwise74/video-api/internal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type partialVerifToken struct {
	ExpiresAt *time.Time
	Purpose   string
	Used      bool
}

func UserVerify(c *gin.Context, d *internal.Deps) {
	requestID := c.MustGet("requestID").(string)

	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No verification token provided",
			"requestID": requestID,
		})
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No user ID provided",
			"requestID": requestID,
		})
		return
	}

	var verifRecord partialVerifToken

	err := d.DB.
		Model(model.VerificationToken{}).
		Where("user_id = ? AND token = ?", userID, token).
		Select("expires_at", "purpose", "used").
		First(&verifRecord).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "Token expired or invalid",
				"requestID": requestID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to get verification token record", zap.Error(err))
		return
	}

	if verifRecord.Used {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Token was used already",
			"requestID": requestID,
		})
		return
	}

	if verifRecord.ExpiresAt != nil && verifRecord.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Token expired",
			"requestID": requestID,
		})
		return
	}

	err = d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.VerificationToken{}).
			Where("user_id = ? AND token = ?", userID, token).
			Updates(map[string]any{
				"used":    true,
				"used_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.User{}).
			Where("id = ?", userID).
			Updates(map[string]any{
				"verified":   true,
				"expires_at": nil,
			}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to validate user",
			"requestID": requestID,
		})
		zap.L().Error("Failed to update user and token in transaction", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "User validated successfully",
		"requestID": requestID,
	})
}
