package file

import (
	"bitwise74/video-api/internal"
	"bitwise74/video-api/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func FileOwns(c *gin.Context, d *internal.Deps) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	zap.L().Debug("FileOwns endpoint called", zap.String("requestID", requestID), zap.String("userID", userID))

	fileID := c.Param("id")
	if fileID == "" {
		zap.L().Debug("No file ID provided", zap.String("requestID", requestID))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No file ID provided",
			"requestID": requestID,
		})
		return
	}

	zap.L().Debug("Checking file ownership", zap.String("fileID", fileID), zap.String("userID", userID))
	var owns bool
	err := d.DB.
		Model(model.File{}).
		Where("id = ? AND user_id = ?", fileID, userID).
		Select("count(*) > 0").
		Find(&owns).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			zap.L().Debug("File not found for ownership check", zap.String("fileID", fileID))
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

	if owns {
		zap.L().Debug("User owns file", zap.String("fileID", fileID), zap.String("userID", userID))
		c.JSON(http.StatusOK, gin.H{"owns": true})
		return
	}

	zap.L().Debug("User does not own file", zap.String("fileID", fileID), zap.String("userID", userID))
	c.JSON(http.StatusForbidden, gin.H{"owns": false})
}
