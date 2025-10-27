package file

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/redis"
	"bitwise74/video-api/internal/types"
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type deleteInfo struct {
	FileKey  string
	ThumbKey string
	Size     int
}

type deleteRequest struct {
	IDs []int `json:"ids" binding:"required"`
}

func Delete(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	var req deleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"requestID": requestID,
		})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No file IDs provided",
			"requestID": requestID,
		})
		return
	}

	var info []deleteInfo

	err := d.DB.Gorm.
		Model(model.File{}).
		Where("user_id = ? AND id IN ?", userID, req.IDs).
		Select("file_key", "size").
		Find(&info).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "File not found. It either doesn't exist or you don't own it",
				"requestID": requestID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to check if file exists", zap.Error(err))
		return
	}

	tx := d.DB.Gorm.Begin()

	err = tx.
		Where("id IN ?", req.IDs).
		Delete(model.File{}).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		tx.Rollback()

		zap.L().Error("Failed to check if file exists", zap.Error(err))
		return
	}

	// Format deleteInfo into something usable
	objects := []awsTypes.ObjectIdentifier{}
	totalSize := 0

	for _, v := range info {
		thumbKey := strings.TrimSuffix(v.FileKey, ".mp4") + ".webp"

		objects = append(objects,
			awsTypes.ObjectIdentifier{Key: &v.FileKey},
			awsTypes.ObjectIdentifier{Key: &thumbKey},
		)

		totalSize += v.Size
	}

	_, err = d.S3.C.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: d.S3.Bucket,
		Delete: &awsTypes.Delete{
			Objects: objects,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		tx.Rollback()

		zap.L().Error("Failed to delete file from S3", zap.Error(err))
		return
	}

	err = tx.
		Model(model.Stats{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"used_storage":   gorm.Expr("used_storage - ?", totalSize),
			"uploaded_files": gorm.Expr("uploaded_files - ?", len(req.IDs)),
		}).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		tx.Rollback()

		zap.L().Error("Failed to decrement user's used storage", zap.Error(err))
		return
	}

	var newStats model.Stats

	err = tx.
		Where("user_id = ?", userID).
		First(&newStats).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		tx.Rollback()

		zap.L().Error("Failed to decrement user's used storage", zap.Error(err))
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		tx.Rollback()

		zap.L().Error("Failed to commit transaction", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, newStats)

	redis.InvalidateCache("user:" + userID)
}
