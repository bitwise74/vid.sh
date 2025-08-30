package ffmpeg

import (
	"bitwise74/video-api/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func FFmpegProgress(c *gin.Context) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	if _, ok := service.ProgressMap.Load(userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "No running jobs found",
			"requestID": requestID,
		})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "nocache")
	c.Header("Connection", "keep-alive")

	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	for range ticker.C {
		val, ok := service.ProgressMap.Load(userID)
		if !ok {
			continue
		}

		v := val.(service.FFMpegJobStats)

		fmt.Fprintf(c.Writer, "data: %.2f\n\n", v.Progress)
		c.Writer.Flush()

		if v.Progress >= 100 {
			break
		}
	}

	fmt.Fprintf(c.Writer, "data: %.2f\n\n", 100.0)
	c.Writer.Flush()
}
