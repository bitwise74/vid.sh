package user

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/types"
	"bitwise74/video-api/pkg/security"
	"bitwise74/video-api/pkg/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type resetPasswdBody struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func ResetPassword(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)

	var body resetPasswdBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Malformed or invalid JSON body",
			"requestID": requestID,
		})

		return
	}

	var userID string
	tx := d.DB.Gorm.Begin()

	err := tx.
		Model(model.Token{}).
		Where("token = ?", body.Token).
		Update("used", true).
		Pluck("user_id", &userID).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Invalid or expired token",
				"requestID": requestID,
			})
			tx.Rollback()
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to mark token as used", zap.Error(err))
		tx.Rollback()
		return
	}

	if err := validators.PasswordValidator(body.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"requestID": requestID,
		})
		tx.Rollback()
		return
	}

	passwdHash, err := security.Argon.GenerateFromPassword(body.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to hash password", zap.Error(err))
		tx.Rollback()
		return
	}

	err = tx.Model(model.User{}).Where("id = ?", userID).Update("password_hash", passwdHash).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to update user's password", zap.Error(err))
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to commit transaction", zap.Error(err))
		tx.Rollback()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Password reset successfully",
		"requestID": requestID,
	})
}
