package ffmpeg

import (
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	jobID := `job-` + userID + `-` + requestID

	service.ProgressMap.Store(userID, service.FFMpegJobStats{
		Progress: 0,
		JobID:    jobID,
		State:    "Waiting in queue...",
		Stopped:  false,
	})

	c.JSON(http.StatusOK, gin.H{
		"jobID": jobID,
	})
}
