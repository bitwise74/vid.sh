package ffmpeg

import (
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/internal/types"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Progress(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	_, ok := service.ProgressMap.Load(userID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No job found for this user",
			"request": requestID,
		})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "nocache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	for range ticker.C {
		// If the job doesn't exist in the map we shouldn't send any more updates
		val, ok := service.ProgressMap.Load(userID)
		if !ok {
			break
		}

		v := val.(service.FFMpegJobStats)

		fmt.Fprintf(c.Writer, "data: %.2f|%s\n\n", v.Progress, v.State)
		c.Writer.Flush()

		if v.Stopped || v.Progress >= 100 {
			break
		}
	}

	fmt.Fprint(c.Writer, "data: 100|Finishing...\n\n")
	c.Writer.Flush()
}
