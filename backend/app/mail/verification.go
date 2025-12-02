package mail

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/redis"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/internal/types"
	"bitwise74/video-api/pkg/security"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VerificationMail(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)
	ctxUser := c.MustGet("ctxUser").(*model.User)

	exp, err := redis.CheckPenalties(c.Request.Context(), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if exp != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       "Too many verification email requests. Please try again later.",
			"retry_after": exp,
			"requestID":   requestID,
		})
		return
	}

	if ctxUser.Verified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already verified"})
		return
	}

	expireAt := time.Now().Add(time.Minute * 30)
	cleanAt := time.Now().Add(time.Hour * 24 * 60)

	token, err := security.MakeToken(&security.TokenOpts{
		UserID:    userID,
		Purpose:   "verification-resend",
		ExpiresAt: &expireAt,
		CleanupAt: &cleanAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to create token", zap.String("requestID", requestID), zap.Error(err))
		return
	}

	tx := d.DB.Gorm.Begin()

	err = tx.Create(token).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to save verification token", zap.String("requestID", requestID), zap.Error(err))
		return
	}

	err = service.SendVerificationMail(token, ctxUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to send verification mail", zap.String("requestID", requestID), zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification mail sent",
	})

	err = redis.AddPenalty(context.Background(), c.ClientIP(), time.Hour*24)
	if err != nil {
		zap.L().Error("Failed to add penalty", zap.String("requestID", requestID), zap.String("ip", c.ClientIP()), zap.Error(err))
	}
}
