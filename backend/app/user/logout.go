package user

import (
	"bitwise74/video-api/internal"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TODO: make this blacklist the JWT until it expires

func Logout(c *gin.Context, d *internal.Deps) {
	sslEnabled, err := strconv.ParseBool("HOST_SSL_ENABLED")
	if err != nil {
		sslEnabled = false
	}

	c.SetCookie("user_id", "", -1, "/", "", sslEnabled, false)
	c.SetCookie("auth_token", "", -1, "/", "", sslEnabled, true)
	c.SetCookie("logged_in", "", -1, "/", "", sslEnabled, false)
}
