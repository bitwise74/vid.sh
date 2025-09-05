package ffmpeg

import (
	"bitwise74/video-api/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FFmpegProgress(c *gin.Context) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	zap.L().Debug("FFmpegProgress request received", zap.String("requestID", requestID), zap.String("userID", userID))

	if _, ok := service.ProgressMap.Load(userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "No running jobs found",
			"requestID": requestID,
		})
		zap.L().Debug("No running jobs found for user", zap.String("requestID", requestID), zap.String("userID", userID))
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "nocache")
	c.Header("Connection", "keep-alive")
	zap.L().Debug("Headers set for SSE stream", zap.String("requestID", requestID))

	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	for range ticker.C {
		val, ok := service.ProgressMap.Load(userID)
		if !ok {
			zap.L().Debug("ProgressMap entry missing during streaming", zap.String("requestID", requestID), zap.String("userID", userID))
			continue
		}

		v := val.(service.FFMpegJobStats)

		fmt.Fprintf(c.Writer, "data: %.2f\n\n", v.Progress)
		c.Writer.Flush()
		zap.L().Debug("Progress sent to client", zap.String("requestID", requestID), zap.String("userID", userID), zap.Float64("progress", v.Progress))

		if v.Progress >= 100 {
			zap.L().Debug("FFmpeg job completed", zap.String("requestID", requestID), zap.String("userID", userID))
			break
		}
	}

	fmt.Fprintf(c.Writer, "data: %.2f\n\n", 100.0)
	c.Writer.Flush()
	zap.L().Debug("Final progress sent to client", zap.String("requestID", requestID), zap.String("userID", userID))
}
