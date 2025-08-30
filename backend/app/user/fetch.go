package user

import (
	"bitwise74/video-api/internal"
	"bitwise74/video-api/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UserFetch(c *gin.Context, d *internal.Deps) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	var videos []model.File

	err := d.DB.
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(10).
		Find(&videos).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to fetch initial user data", zap.Error(err))
		return
	}

	for i, file := range videos {
		version := strconv.Itoa(file.Version)
		videos[i].FileKey = file.FileKey + "?v=" + version
		videos[i].ThumbKey = file.ThumbKey + "?v=" + version
	}

	var stats model.Stats
	err = d.DB.
		Where("user_id = ?", userID).
		First(&stats).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to fetch initial user data", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
		"stats":  stats,
	})
}
