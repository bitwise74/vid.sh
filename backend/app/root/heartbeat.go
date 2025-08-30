package root

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Heartbeat(c *gin.Context) {
	c.Status(http.StatusOK)
}
