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

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No file ID provided",
			"requestID": requestID,
		})
		return
	}

	var owns bool
	err := d.DB.
		Model(model.File{}).
		Where("id = ? AND user_id = ?", fileID, userID).
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

	if owns {
		c.JSON(http.StatusOK, gin.H{"owns": true})
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"owns": false})
}
