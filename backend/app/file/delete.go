package file

import (
	"bitwise74/video-api/internal"
	"bitwise74/video-api/internal/model"
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type deleteInfo struct {
	FileKey  string
	ThumbKey string
	Size     int
}

func FileDelete(c *gin.Context, d *internal.Deps) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "ID is missing",
			"requestID": requestID,
		})
		return
	}

	var info deleteInfo

	err := d.DB.
		Model(model.File{}).
		Where("user_id = ? AND id = ?", userID, fileID).
		Select("file_key", "thumb_key", "size").
		First(&info).
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

	err = d.DB.
		Where("file_key = ?", info.FileKey).
		Delete(model.File{}).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to check if file exists", zap.Error(err))
		return
	}

	resp, err := d.S3.C.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: d.S3.Bucket,
		Delete: &types.Delete{
			Objects: []types.ObjectIdentifier{
				{Key: &info.FileKey},
				{Key: &info.ThumbKey},
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to delete file from S3", zap.Error(err))
		return
	}

	for _, v := range resp.Deleted {
		zap.L().Debug("Deleted item", zap.String("item", *v.Key))
	}

	err = d.DB.
		Model(model.Stats{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"used_storage":   gorm.Expr("used_storage - ?", info.Size),
			"uploaded_files": gorm.Expr("uploaded_files - ?", 1),
		}).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to decrement user's used storage", zap.Error(err))
		return
	}

	var newStats model.Stats

	err = d.DB.
		Where("user_id = ?", userID).
		First(&newStats).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to decrement user's used storage", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, newStats)
}
