package root

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Validate(c *gin.Context) {
	c.Status(http.StatusOK)
}
