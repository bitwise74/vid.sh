package user

import (
	"bitwise74/video-api/internal/types"
	"bitwise74/video-api/pkg/security"
	"bitwise74/video-api/pkg/validators"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

func Login(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)

	var data loginBody
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"requestID": requestID,
		})

		zap.L().Error("Can't bind request body",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return
	}

	if err := validators.EmailValidator(data.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
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

	user, err := d.DB.FetchUserByEmail(data.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":     "Invalid credentials",
				"requestID": requestID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to fetch user by email",
			zap.String("requestID", requestID),
			zap.String("email", data.Email),
			zap.Error(err),
		)
		return
	}

	ok, err := security.Argon.VerifyPasswd(data.Password, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to verify password",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
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

		zap.L().Error("Failed to generate JWT auth token",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return
	}

	sslEnabled, err := strconv.ParseBool(os.Getenv("HOST_SSL_ENABLED"))
	if err != nil {
		sslEnabled = false
	}

	if !data.Remember {
		c.SetCookie("logged_in", "1", 0, "/", "", sslEnabled, false)
		c.SetCookie("auth_token", authToken, 0, "/", "", sslEnabled, true)
	} else {
		c.SetCookie("auth_token", authToken, 60*60*24*30, "/", "", sslEnabled, true)
		c.SetCookie("logged_in", "1", 60*60*24*30, "/", "", sslEnabled, false)
	}

	c.JSON(http.StatusOK, gin.H{
		"userID":   user.ID,
		"verified": user.Verified,
	})
}

func makeToken(c *jwt.MapClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString([]byte(os.Getenv("SECURITY_JWT_SECRET")))
}
