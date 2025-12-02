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
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TODO: use old email
type ChangeEmailBody struct {
	NewEmail string `json:"new_email"`
}

func ChangeEmail(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)
	oldEmail := c.MustGet("ctxUser").(*model.User).Email

	// Check for penalties
	exp, err := redis.CheckPenalties(c.Request.Context(), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to check penalties",
			zap.String("requestID", requestID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return
	}

	if exp != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       "Too many email change attempts. Please try again later.",
			"retry_after": exp,
			"requestID":   requestID,
		})
		return
	}

	var body ChangeEmailBody
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"requestID": requestID,
		})
		return
	}

	expiresAt := time.Now().Add(time.Minute * 30)
	cleanupAt := time.Now().Add(time.Hour * 24 * 40)

	token, err := security.MakeToken(&security.TokenOpts{
		UserID:    userID,
		Purpose:   "change-email",
		ExpiresAt: &expiresAt,
		CleanupAt: &cleanupAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to create token",
			zap.String("requestID", requestID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return
	}

	err = d.DB.SaveToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to save token to db",
			zap.String("requestID", requestID),
			zap.String("userID", userID),
			zap.Error(err),
		)
	}

	err = service.SendMail(oldEmail, "Change your email for vid.sh", fmt.Sprintf("Click <a href=\"%s/change_email?user_id=%s&new_email=%s\">here</a> to confirm your new email address. If you didn't request this you can safely ignore this email", c.Request.Host, userID, body.NewEmail))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to send confirmation email",
			"requestID": requestID,
		})

		zap.L().Error("Failed to send change email confirmation",
			zap.String("requestID", requestID),
			zap.String("userID", userID),
			zap.String("newEmail", body.NewEmail),
			zap.Error(err),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Confirmation email sent",
	})

	err = redis.AddPenalty(context.Background(), c.ClientIP(), time.Hour*24)
	if err != nil {
		zap.L().Error("Failed to add penalty", zap.String("requestID", requestID), zap.String("ip", c.ClientIP()), zap.Error(err))
	}
}
