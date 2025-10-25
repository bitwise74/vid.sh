package user

import (
	"bitwise74/video-api/internal/redis"
	"bitwise74/video-api/internal/types"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Fetch(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	if redis.CheckCache("user:"+userID, c) {
		return
	}

	user, err := d.DB.FetchUserDataByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to fetch user data",
			zap.String("requestID", requestID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return
	}

	c.JSON(http.StatusOK, user)

	if err := redis.Rdb.Set(context.Background(), "user:"+userID, user, time.Minute*5).Err(); err != nil {
		zap.L().Error("Failed to set user cache",
			zap.String("requestID", requestID),
			zap.String("userID", userID),
			zap.Error(err),
		)
	}
}
