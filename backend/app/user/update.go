package user

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/redis"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/internal/types"
	"bitwise74/video-api/pkg/util"
	"bitwise74/video-api/pkg/validators"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Update(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	var data validators.UserUpdateOpts
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Malformed or invalid multipart form",
			"requestID": requestID,
		})
		return
	}

	if data.IsEmpty() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No update options provided",
			"requestID": requestID,
		})
		return
	}

	if err := validators.UserValidator(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"requestID": requestID,
		})
		return
	}

	taken, err := d.DB.CheckIfUsernameTaken(data.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to check if username is taken", zap.String("requestID", requestID), zap.Error(err))
		return
	}

	if taken {
		c.JSON(http.StatusConflict, gin.H{
			"error":     "A user with that username already exists",
			"requestID": requestID,
		})
		return
	}

	user, err := d.DB.FetchUserDataByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to load user from database", zap.String("requestID", requestID), zap.Error(err))
		return
	}

	if data.Username != "" {
		user.Username = data.Username
	}

	if data.Avatar != nil {
		file, err := os.CreateTemp("", "profile-pic-og-*")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to create temporary file for profile picture", zap.String("requestID", requestID), zap.Error(err))
			return
		}
		defer file.Close()
		defer os.Remove(file.Name())

		fh, err := data.Avatar.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to read uploaded file", zap.String("requestID", requestID), zap.Error(err))
			return
		}

		_, err = io.Copy(file, fh)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to copy file contents to destination", zap.String("requestID", requestID), zap.Error(err))
			return
		}

		if data.AvatarCrop != "" {
			path, err := service.Crop(file.Name(), data.AvatarCropProper)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":     "Internal server error",
					"requestID": requestID,
				})

				zap.L().Error("Failed to process profile picture", zap.String("requestID", requestID), zap.Error(err))

				if path != "" {
					os.Remove(path)
				}

				return
			}
			defer os.Remove(path)

			file, err = os.Open(path)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":     "Internal server error",
					"requestID": requestID,
				})

				zap.L().Error("Failed to open processed profile picture file", zap.String("requestID", requestID), zap.Error(err))
				return
			}
			// No defer because reassignment
		}

		key := util.RandStr(32) + ".webp"

		file.Seek(0, 0)

		if user.AvatarHash != "" {
			_, err = d.S3.C.DeleteObject(context.Background(), &s3.DeleteObjectInput{
				Bucket: d.S3.Bucket,
				Key:    aws.String("avatars/" + user.AvatarHash),
			})
			if err != nil {
				zap.L().Error("Failed to delete old profile picture", zap.String("requestID", requestID), zap.Error(err))
			}
		}

		// Upload to the /avatars bucket directory
		// TODO: add cleanup on errors
		_, err = d.S3.C.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket:       d.S3.Bucket,
			Key:          aws.String("avatars/" + key),
			Body:         file,
			CacheControl: aws.String("public, max-age=14400"),
			ContentType:  aws.String("image/webp"),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to upload profile picture", zap.String("requestID", requestID), zap.Error(err))
			return
		}

		user.AvatarHash = key
	}

	if data.PublicProfileEnabled != nil {
		user.PublicProfileEnabled = *data.PublicProfileEnabled
	}

	err = d.DB.Gorm.
		Model(model.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"avatar_hash":            gorm.Expr("?", user.AvatarHash),
			"username":               gorm.Expr("?", user.Username),
			"public_profile_enabled": gorm.Expr("?", user.PublicProfileEnabled),
		}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed update user profile", zap.String("requestID", requestID), zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":   user.Username,
		"avatarHash": user.AvatarHash,
	})

	redis.InvalidateCache("user:" + userID)
}
