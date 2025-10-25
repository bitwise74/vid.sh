package user

import (
	"bitwise74/video-api/internal/types"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func Logout(c *gin.Context, d *types.Dependencies) {
	tokenStr, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No auth_token cookie",
			"requestID": c.MustGet("requestID").(string),
		})
		return
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}

		return []byte("SECURITY_JWT_SECRET"), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":     "Authorization token invalid",
			"requestID": c.MustGet("requestID").(string),
		})
		return
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":     "Authorization token invalid",
			"requestID": c.MustGet("requestID").(string),
		})
		return
	}

	expiry := time.Unix(int64(expFloat), 0)
	err = d.DB.BlacklistJWTToken(tokenStr, expiry)
	if err != nil {
		zap.L().Error("Failed to blacklist JWT token", zap.Error(err))
	}

	sslEnabled, err := strconv.ParseBool("HOST_SSL_ENABLED")
	if err != nil {
		sslEnabled = false
	}

	c.SetCookie("auth_token", "", -1, "/", "", sslEnabled, true)
	c.SetCookie("logged_in", "", -1, "/", "", sslEnabled, false)
}
