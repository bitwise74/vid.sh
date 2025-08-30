// Package middleware contains any custom middleware used in the app
package middleware

import (
	"bitwise74/video-api/pkg/util"

	"github.com/gin-gonic/gin"
)

// NewRequestIDMiddleware returns a new middleware function that generates a request ID for
// each incoming request and sets it as requestID
func NewRequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("requestID", util.RandStr(10))
		c.Next()
	}
}
