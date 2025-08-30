package ffmpeg

import (
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FFMpegStart(c *gin.Context) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	if _, ok := service.ProgressMap.Load(userID); ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error":     "A job is running already. Wait for it to finish first",
			"requestID": requestID,
		})
		return
	}

	jobID := util.RandStr(5)
	service.ProgressMap.Store(userID, service.FFMpegJobStats{
		Progress: 0.0,
		JobID:    jobID,
	})

	c.JSON(http.StatusOK, gin.H{
		"jobID": jobID,
	})
}
