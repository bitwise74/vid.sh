package ffmpeg

import (
	"bitwise74/video-api/internal"
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/service"
	"bitwise74/video-api/pkg/util"
	"bitwise74/video-api/pkg/validators"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func FFmpegProcess(c *gin.Context, d *internal.Deps) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)
	jobID := c.Query("jobID")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No job ID provided",
			"requestID": requestID,
		})
		return
	}

	var opts validators.ProcessingOptions
	if err := c.MustBindWith(&opts, binding.FormMultipart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to read form body",
			"requestID": requestID,
		})

		zap.L().Error("Failed to read form body", zap.Error(err))
		return
	}

	if code, err := validators.ProcessingOptsValidator(&opts, float64(opts.File.Size)); err != nil {
		c.JSON(code, gin.H{
			"error":     err.Error(),
			"requestID": requestID,
		})
		return
	}

	code, f, err := validators.FileValidator(opts.File, nil, "")
	if err != nil {
		if code == http.StatusInternalServerError {
			zap.L().Error("Failed to validate file", zap.Error(err))

			err = errors.New("Internal server error")
		}

		c.JSON(code, gin.H{
			"error":     err.Error(),
			"requestID": requestID,
		})
		return
	}
	defer f.Close()

	tempFile, err := os.CreateTemp("", "upload-*.mp4")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to create temporary file", zap.Error(err))
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to copy data to temporary file", zap.Error(err))
		return
	}

	if !opts.SaveToCloud {
		c.Header("Content-Type", "video/mp4")
		c.Header("Transfer-Encoding", "chunked")

		ctxReq := c.Request.Context()
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		ctx, cancelMerged := util.MergeContexts(ctxReq, ctxTimeout)
		defer cancelMerged()

		done := make(chan error, 1)
		// Enqueue can only error if the queue is full
		err = d.JobQueue.Enqueue(&service.FFmpegJob{
			ID:       jobID,
			UserID:   userID,
			FilePath: tempFile.Name(),
			Output:   c.Writer,
			Opts:     &opts,
			UseGPU:   true,
			Ctx:      ctx,
			Done:     done,
		})
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":     "Job queue is full. Please wait a moment before trying again",
				"requestID": requestID,
			})

			zap.L().Warn("FFmpeg job queue is full")
			return
		}

		select {
		case err := <-done:
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":     "Internal server error",
					"requestID": requestID,
				})
				return
			}
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":     "Request was cancelled or timed out",
				"requestID": requestID,
			})

			zap.L().Warn("Request context done before FFmpeg finished", zap.Error(ctx.Err()))
			return
		}
		return
	}

	tempProcessed, err := os.CreateTemp("", "processed-*.mp4")
	if err != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Warn("Failed to create temp file for processing", zap.Error(err))
		return
	}

	ctxReq := c.Request.Context()
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ctx, cancelMerged := util.MergeContexts(ctxReq, ctxTimeout)
	defer cancelMerged()

	done := make(chan error, 1)
	// Enqueue can only error if the queue is full
	err = d.JobQueue.Enqueue(&service.FFmpegJob{
		ID:       jobID,
		UserID:   userID,
		FilePath: tempFile.Name(),
		Output:   tempProcessed,
		Opts:     &opts,
		UseGPU:   true,
		Ctx:      ctx,
		Done:     done,
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":     "Job queue is full. Please wait a moment before trying again",
			"requestID": requestID,
		})

		zap.L().Warn("FFmpeg job queue is full")
		return
	}

	select {
	case err := <-done:
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})
			return
		}
	case <-ctx.Done():
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error":     "Request was cancelled or timed out",
			"requestID": requestID,
		})

		zap.L().Warn("Request context done before FFmpeg finished", zap.Error(ctx.Err()))
		return
	}

	fileEnt, err := d.Uploader.Do(tempProcessed.Name(), opts.File.Filename, userID)
	if err != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Warn("Failed to upload file to S3", zap.Error(err))
		return
	}

	err = d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&fileEnt).Error; err != nil {
			return err
		}

		if err := tx.
			Model(model.Stats{}).
			Where("user_id = ?", userID).
			Updates(map[string]any{
				"used_storage":   gorm.Expr("used_storage + ?", fileEnt.Size),
				"uploaded_files": gorm.Expr("uploaded_files + ?", 1),
			}).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Database transaction failed", zap.String("requestID", requestID), zap.Error(err))
		return
	}

	c.Status(http.StatusOK)
}
