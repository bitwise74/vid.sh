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

	zap.L().Debug("FFmpeg request received", zap.String("requestID", requestID), zap.String("userID", userID), zap.String("jobID", jobID))

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
		zap.L().Error("Failed to read form body", zap.Error(err), zap.String("requestID", requestID))
		return
	}
	zap.L().Debug("Form body parsed successfully", zap.String("requestID", requestID), zap.String("filename", opts.File.Filename))

	if code, err := validators.ProcessingOptsValidator(&opts, float64(opts.File.Size)); err != nil {
		c.JSON(code, gin.H{
			"error":     err.Error(),
			"requestID": requestID,
		})
		zap.L().Debug("Processing options validation failed", zap.String("requestID", requestID), zap.Error(err))
		return
	}
	zap.L().Debug("Processing options validated", zap.String("requestID", requestID))

	code, f, err := validators.FileValidator(opts.File, nil, "")
	if err != nil {
		if code == http.StatusInternalServerError {
			zap.L().Error("Failed to validate file", zap.Error(err), zap.String("requestID", requestID))
			err = errors.New("Internal server error")
		}
		c.JSON(code, gin.H{
			"error":     err.Error(),
			"requestID": requestID,
		})
		return
	}
	defer f.Close()
	zap.L().Debug("File validated successfully", zap.String("requestID", requestID), zap.String("filename", opts.File.Filename))

	tempFile, err := os.CreateTemp("", "upload-*.mp4")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		zap.L().Error("Failed to create temporary file", zap.Error(err), zap.String("requestID", requestID))
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())
	zap.L().Debug("Temporary file created", zap.String("requestID", requestID), zap.String("tempFile", tempFile.Name()))

	_, err = io.Copy(tempFile, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		zap.L().Error("Failed to copy file to temp file", zap.Error(err), zap.String("requestID", requestID))
		return
	}
	zap.L().Debug("File copied to temporary file", zap.String("requestID", requestID))

	if !opts.SaveToCloud {
		c.Header("Content-Type", "video/mp4")
		c.Header("Transfer-Encoding", "chunked")

		ctxReq := c.Request.Context()
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		ctx, cancelMerged := util.MergeContexts(ctxReq, ctxTimeout)
		defer cancelMerged()

		done := make(chan error, 1)
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
			zap.L().Warn("FFmpeg job queue is full", zap.String("requestID", requestID))
			return
		}
		zap.L().Debug("FFmpeg job enqueued", zap.String("requestID", requestID))

		select {
		case err := <-done:
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":     "Internal server error",
					"requestID": requestID,
				})
				zap.L().Error("FFmpeg job failed", zap.String("requestID", requestID), zap.Error(err))
				return
			}
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":     "Request was cancelled or timed out",
				"requestID": requestID,
			})
			zap.L().Warn("Request context done before FFmpeg finished", zap.String("requestID", requestID), zap.Error(ctx.Err()))
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
		zap.L().Warn("Failed to create temp file for processed video", zap.String("requestID", requestID), zap.Error(err))
		return
	}

	ctxReq := c.Request.Context()
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	ctx, cancelMerged := util.MergeContexts(ctxReq, ctxTimeout)
	defer cancelMerged()

	done := make(chan error, 1)
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
		zap.L().Warn("FFmpeg job queue is full (cloud save)", zap.String("requestID", requestID))
		return
	}
	zap.L().Debug("FFmpeg job enqueued for cloud save", zap.String("requestID", requestID))

	select {
	case err := <-done:
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal server error",
				"requestID": requestID,
			})
			zap.L().Error("FFmpeg job failed (cloud save)", zap.String("requestID", requestID), zap.Error(err))
			return
		}
	case <-ctx.Done():
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error":     "Request was cancelled or timed out",
			"requestID": requestID,
		})
		zap.L().Warn("Request context done before FFmpeg finished (cloud save)", zap.String("requestID", requestID), zap.Error(ctx.Err()))
		return
	}

	fileEnt, err := d.Uploader.Do(tempProcessed.Name(), opts.File.Filename, userID)
	if err != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		zap.L().Warn("Failed to upload file to S3", zap.String("requestID", requestID), zap.Error(err))
		return
	}
	zap.L().Debug("File uploaded to cloud storage", zap.String("requestID", requestID), zap.String("filename", opts.File.Filename))

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
	zap.L().Debug("Database transaction completed successfully", zap.String("requestID", requestID))

	c.Status(http.StatusOK)
	zap.L().Debug("FFmpeg request finished successfully", zap.String("requestID", requestID))
}
