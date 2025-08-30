package user

import (
	"bitwise74/video-api/internal"
	"bitwise74/video-api/internal/model"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserLogin(c *gin.Context, d *internal.Deps) {
	requestID := c.MustGet("requestID").(string)

	var data loginBody
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"requestID": requestID,
		})

		zap.L().Error("Can't bind request body", zap.Error(err), zap.String("requestID", requestID))
		return
	}

	if data.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Email field can't be empty",
			"requestID": requestID,
		})
		return
	}

	if data.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Password field can't be empty",
			"requestID": requestID,
		})
		return
	}

	var user model.User

	if err := d.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "User not found",
			"requestID": requestID,
		})

		zap.L().Error("User not found", zap.Error(err), zap.String("requestID", requestID))
		return
	}

	ok, err := d.Argon.VerifyPasswd(data.Password, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to verify password", zap.Error(err), zap.String("requestID", requestID))
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":     "Invalid credentials",
			"requestID": requestID,
		})
		return
	}

	authToken, err := makeToken(&jwt.MapClaims{
		"user_id": user.ID,
		"type":    "auth",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to generate JWT auth token", zap.Error(err), zap.String("requestID", requestID))
		return
	}

	sslEnabled, err := strconv.ParseBool("HOST_SSL_ENABLED")
	if err != nil {
		sslEnabled = false
	}

	c.SetCookie("user_id", user.ID, 9999999, "/", "", sslEnabled, false)
	c.SetCookie("auth_token", authToken, 60*60*24*30, "/", "", sslEnabled, true)
	c.SetCookie("logged_in", "1", 60*60*24*30, "/", "", sslEnabled, false)
	c.JSON(http.StatusOK, gin.H{
		"userID":   user.ID,
		"verified": user.Verified,
	})
}

func makeToken(c *jwt.MapClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString([]byte(os.Getenv("SECURITY_JWT_SECRET")))
}
