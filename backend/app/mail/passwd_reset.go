package mail

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/redis"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/internal/types"
	"bitwise74/video-api/pkg/security"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PasswdResetBody struct {
	Email string `json:"email"`
}

type partialUser struct {
	Verified bool
	Email    string
	ID       string
}

func PasswdReset(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)

	var body PasswdResetBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Can't bind request body",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return
	}

	if body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Email is required",
			"requestID": requestID,
		})
		return
	}

	// Check for penalties
	val := redis.Rdb.Get(c.Request.Context(), c.ClientIP()+":penalty").Val()
	if val != "" {
		valInt, _ := strconv.Atoi(val)

		if valInt >= 3 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":     "Too many requests",
				"requestID": requestID,
			})
		}
	}

	var user partialUser
	err := d.DB.Gorm.
		Model(model.User{}).
		Select("verified", "email", "id").
		Where("email = ?", body.Email).
		First(&user).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "User not found",
				"requestID": requestID,
			})
			return
		}

		zap.L().Error("Error fetching user", zap.String("requestID", requestID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Verify your email address first",
			"requestID": requestID,
		})
		return
	}

	expiresAt := time.Now().Add(30 * time.Minute)
	cleanupAt := time.Now().Add(24 * 60 * time.Hour)

	resetToken, err := security.MakeToken(&security.TokenOpts{
		UserID:    user.ID,
		Purpose:   "password-reset",
		ExpiresAt: &expiresAt,
		CleanupAt: &cleanupAt,
	})
	if err != nil {
		zap.L().Error("Failed to create password reset token",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		return
	}

	err = d.DB.Gorm.Create(&resetToken).Error
	if err != nil {
		zap.L().Error("Failed to save password reset token",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		return
	}

	// Try to send mail now
	err = service.SendMail(user.Email, "Reset your password for vid.sh", fmt.Sprintf("Click <a href='https://%v/reset-passwd?token=%v'>here</a> to reset your password.\n\nThis link will expire in 30 minutes", c.Request.Host, resetToken.Token))
	if err != nil {
		zap.L().Error("Failed to send password reset email",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent",
	})

	err = redis.AddPenalty(context.Background(), c.ClientIP(), time.Hour*24)
	if err != nil {
		zap.L().Error("Failed to add penalty", zap.String("requestID", requestID), zap.String("email", body.Email), zap.Error(err))
	}
}
