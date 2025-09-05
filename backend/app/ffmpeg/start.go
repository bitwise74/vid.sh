package ffmpeg

import (
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FFMpegStart(c *gin.Context) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	zap.L().Debug("FFMpegStart request received", zap.String("requestID", requestID), zap.String("userID", userID))

	if _, ok := service.ProgressMap.Load(userID); ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error":     "A job is running already. Wait for it to finish first",
			"requestID": requestID,
		})
		zap.L().Debug("Job already running for user", zap.String("requestID", requestID), zap.String("userID", userID))
		return
	}

	jobID := util.RandStr(5)
	service.ProgressMap.Store(userID, service.FFMpegJobStats{
		Progress: 0.0,
		JobID:    jobID,
	})
	zap.L().Debug("FFMpeg job started", zap.String("requestID", requestID), zap.String("userID", userID), zap.String("jobID", jobID))

	c.JSON(http.StatusOK, gin.H{
		"jobID": jobID,
	})
}
