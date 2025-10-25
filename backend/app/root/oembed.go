package root

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OEmbed(c *gin.Context) {
	title := c.Query("title")
	username := c.Query("username")

	c.JSON(http.StatusOK, gin.H{
		"verison":       "1.0",
		"author_name":   title,
		"provider_name": "Video by @" + username,
	})
}
