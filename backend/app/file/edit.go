package file

import (
	"bitwise74/video-api/internal"
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/pkg/util"
	"bitwise74/video-api/pkg/validators"
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type fileEditOpts struct {
	NewName           *string                       `json:"name,omitempty"`
	ProcessingOptions *validators.ProcessingOptions `json:"processing_options,omitempty"`
}

func FileEdit(c *gin.Context, d *internal.Deps) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No file ID provided",
			"requestID": requestID,
		})

		return
	}

	var data fileEditOpts
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Malformed or invalid JSON request body",
			"requestID": requestID,
		})

		zap.L().Error("Failed to read JSON body", zap.Error(err))
		return
	}

	if data.NewName == nil && data.ProcessingOptions == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No edit options provided",
			"requestID": requestID,
		})
		return
	}

	if data.NewName != nil && *data.NewName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Empty name",
			"requestID": requestID,
		})
		return
	}

	var file model.File
	err := d.DB.
		Where("user_id = ? AND id = ?", userID, fileID).
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

	if data.NewName != nil {
		file.OriginalName = *data.NewName
	}

	originalSize := file.Size

	if data.ProcessingOptions != nil {
		if code, err := validators.ProcessingOptsValidator(data.ProcessingOptions, float64(file.Size)); err != nil {
			c.JSON(code, gin.H{
				"error":     err.Error(),
				"requestID": requestID,
			})
			return
		}

		// Download the video to process
		temp, err := os.CreateTemp("", "process-*.mp4")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to create temporary file", zap.Error(err))
			return
		}
		defer temp.Close()
		defer os.Remove(temp.Name())

		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, os.Getenv("CLOUDFRONT_URL")+"/"+file.FileKey, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to prepare download request", zap.Error(err))
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to download video from cloudfront", zap.Error(err))
			return
		}

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{
				"error":     "Failed to fetch file",
				"requestID": requestID,
			})
			zap.L().Error("CloudFront returned non-OK status", zap.Int("status", resp.StatusCode))
			return
		}
		defer resp.Body.Close()

		if _, err := io.Copy(temp, resp.Body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to copy file to temp", zap.Error(err))
			return
		}
		temp.Seek(0, 0)

		ctxReq := c.Request.Context()
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		ctx, cancelMerged := util.MergeContexts(ctxReq, ctxTimeout)
		defer cancelMerged()

		done := make(chan error, 1)

		tempProcessed, err := os.CreateTemp("", "processed-*.mp4")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to create processed file", zap.Error(err))
			return
		}
		defer tempProcessed.Close()
		defer os.Remove(tempProcessed.Name())

		err = d.JobQueue.Enqueue(&service.FFmpegJob{
			ID:       util.RandStr(5),
			UserID:   userID,
			FilePath: temp.Name(),
			Opts:     data.ProcessingOptions,
			UseGPU:   true,
			Output:   tempProcessed,
			Ctx:      ctx,
			Done:     done,
		})
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":     "FFmpeg job queue is full. Please try again later",
				"requestID": requestID,
			})
			return
		}
		if err := <-done; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("FFmpeg failed", zap.Error(err))
			return
		}

		keyNoExt := strings.TrimSuffix(file.FileKey, path.Ext(file.FileKey))

		newFile, err := d.Uploader.Do(tempProcessed.Name(), file.OriginalName, userID, keyNoExt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})

			zap.L().Error("Failed to upload edited video to S3", zap.Error(err))
			return
		}

		zap.L().Debug("New object put to s3")

		file.Duration = newFile.Duration
		file.Size = newFile.Size
	}

	file.Version++

	err = d.DB.Transaction(
		func(tx *gorm.DB) error {
			err := tx.Updates(file).Error
			if err != nil {
				return err
			}

			if originalSize != file.Size {
				err := tx.
					Model(model.Stats{}).
					Where("user_id = ?", userID).
					Updates(map[string]any{
						"used_storage": gorm.Expr("used_storage - ?", originalSize-file.Size),
					}).Error
				if err != nil {
					return err
				}
			}

			return nil
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to commit transaction after file edit", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, file)
}
