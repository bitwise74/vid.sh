package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewBodySizeLimiter(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fast reject for legit requests
		if c.Request.ContentLength > maxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request body size exceeds limit",
			})
			c.Abort()
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
			c.Abort()
			return
		}
	}
}
