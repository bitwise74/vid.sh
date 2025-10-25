package file

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type partialUserData struct {
	AvatarHash string `json:"avatarHash"`
	Username   string `json:"username"`
}

func Fetch(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No file ID provided",
			"requestID": requestID,
		})

		return
	}

	var file model.File

	err := d.DB.Gorm.
		Where("id = ? AND private = ?", fileID, false).
		First(&file).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "File not found",
				"requestID": requestID,
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to fetch file from db", zap.Error(err))
		return
	}

	var user partialUserData
	err = d.DB.Gorm.
		Model(model.User{}).
		Where("id = ?", file.UserID).
		Select("avatar_hash", "username").
		First(&user).
		Error
	if err != nil {
		zap.L().Error("Failed to fetch partial user data. Continuing...", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"file": file,
		"user": user,
	})
}
