package file

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/redis"
	"bitwise74/video-api/internal/types"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func FileOwns(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	if redis.CheckCache("file_owns:"+userID+":"+c.Param("id"), c) {
		return
	}

	fileKey := c.Param("id")
	if fileKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No file ID provided",
			"requestID": requestID,
		})
		return
	}

	var owns bool
	err := d.DB.Gorm.
		Model(model.File{}).
		Where("file_key = ? AND user_id = ?", fileKey, userID).
		Select("count(*) > 0").
		Find(&owns).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "File not found",
				"requestID": requestID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to check if user owns a file", zap.Error(err))
		return
	}

	code := http.StatusOK
	if !owns {
		code = http.StatusForbidden
	}

	c.JSON(code, gin.H{"owns": owns})

	if err := redis.Rdb.Set(context.TODO(), "file_owns:"+userID+":"+fileKey, `{"owns": `+strconv.FormatBool(owns)+`}`, time.Minute*15).Err(); err != nil {
		zap.L().Error("Failed to set cache for file owns", zap.Error(err))
	}
}
