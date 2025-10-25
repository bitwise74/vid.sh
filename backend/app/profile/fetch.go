package profile

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/redis"
	"bitwise74/video-api/internal/types"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type partialProfile struct {
	ID                   string
	AvatarHash           string
	Username             string
	PublicProfileEnabled bool
}

type partialVideo struct {
	FileKey      string  `json:"file_key"`
	OriginalName string  `json:"name"`
	Duration     float64 `json:"duration"`
	CreatedAt    int64   `json:"created_at"`
}

type ProfileResponse struct {
	Username   string         `json:"username"`
	AvatarHash string         `json:"avatarHash"`
	Public     bool           `json:"public"`
	Videos     []partialVideo `json:"videos"`
}

func Fetch(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)

	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No username provided",
		})
		return
	}

	if redis.CheckCache("profile:"+username, c) {
		return
	}

	var prof partialProfile
	err := d.DB.Gorm.
		Model(model.User{}).
		Select("avatar_hash", "username", "id", "public_profile_enabled").
		Where("username = ?", username).
		First(&prof).
		Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Profile not found",
			"requestID": requestID,
		})
		return
	}

	var videos []partialVideo

	if prof.PublicProfileEnabled {
		err = d.DB.Gorm.
			Model(model.File{}).
			Where("private = ? AND user_id = ?", false, prof.ID).
			Select("file_key", "original_name", "duration", "created_at").
			Limit(25).
			Find(&videos).
			Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "Profile not found",
				"requestID": requestID,
			})
			return
		}
	}

	data := ProfileResponse{
		Username:   prof.Username,
		AvatarHash: prof.AvatarHash,
		Public:     prof.PublicProfileEnabled,
		Videos:     videos,
	}

	c.JSON(http.StatusOK, data)

	if err := redis.Rdb.Set(context.Background(), "profile:"+username, data, time.Minute*5).Err(); err != nil {
		zap.L().Error("Failed to set cache for profile fetch", zap.Error(err))
	}
}
