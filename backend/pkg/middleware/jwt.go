package middleware

import (
	"bitwise74/video-api/internal/model"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewJWTMiddleware(d *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.MustGet("requestID").(string)

		tokenStr, err := c.Cookie("auth_token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":     "No auth_token cookie",
					"requestID": requestID,
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":     "Please verify your account before using the service",
				"requestID": requestID,
			})

			zap.L().Error("Failed to get token cookie", zap.Error(err))
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
			}

			return []byte(os.Getenv("SECURITY_JWT_SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":     "Authorization token invalid",
				"requestID": requestID,
			})

			zap.L().Error("Failed to parse token", zap.Error(err), zap.String("requestID", requestID))
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token invalid",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":     "Authorization token invalid",
				"requestID": requestID,
			})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})
			return
		}

		expRaw, ok := claims["exp"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":     "Authorization token expired. Please log in again",
				"requestID": requestID,
			})
			return
		}

		exp, ok := expRaw.(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})
			return
		}

		if time.Now().Unix() >= int64(exp) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":     "Authorization token expired. Please log in again",
				"requestID": requestID,
			})
			return
		}

		// In case someone logs in to delete their account and then logs in again (that's so fucking stupid I know),
		// we'll reject the request
		var user model.User
		err = d.Where("id = ?", userID).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error":     "User not found",
					"requestID": requestID,
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to check if user exists", zap.Error(err), zap.String("requestID", requestID))
			return
		}

		if !user.Verified {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":     "Please verify your account before using the service",
				"requestID": requestID,
			})

			sslEnabled, err := strconv.ParseBool("HOST_SSL_ENABLED")
			if err != nil {
				sslEnabled = false
			}

			// ?????????? what are you doing dumbass
			// 30 days to verify or account deleted
			c.SetCookie("needs_verification", "1", 86400, "/", "", sslEnabled, false)
			return
		}

		// Write some things to the context for later use
		c.Set("userDefaultPrivateVideos", user.DefaultPrivateVideos)
		c.Set("userID", userID)

		c.Next()
	}
}
