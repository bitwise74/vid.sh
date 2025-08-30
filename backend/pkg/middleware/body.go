package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BodySizeLimiter(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fast reject for legit requests
		if c.Request.ContentLength > maxBytes {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Request body size exceeds limit",
			})
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()

		if c.Errors.Last() != nil {
			if strings.Contains(c.Errors.Last().Error(), "http: request body too large") {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "Request body size exceeds limit",
				})
			}
			return
		}
	}
}
