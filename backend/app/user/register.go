package user

import (
	"bitwise74/video-api/db"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/internal/types"
	"bitwise74/video-api/pkg/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type registerBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)

	var data registerBody
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Can't bind request body",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return
	}

	if err := validators.EmailValidator(data.Email); err != nil {
		zap.L().Debug("Invalid email",
			zap.String("requestID", requestID),
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"requestID": requestID,
		})
		return
	}

	if err := validators.PasswordValidator(data.Password); err != nil {
		zap.L().Debug("Invalid password",
			zap.String("requestID", requestID),
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"requestID": requestID,
		})
		return
	}

	user, err := d.DB.CreateUserWithToken(data.Email, data.Password)
	if err != nil {
		if err == db.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{
				"error":     "User with this email already exists",
				"requestID": requestID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to create user",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return
	}

	err = service.SendVerificationMail(&user.Tokens[0], data.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to send verification email",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account created. Please check your email to verify your account.",
	})
}
