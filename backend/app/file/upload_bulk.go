package file

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/types"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UploadFileBulk(c *gin.Context, d *types.Dependencies) {
	c.Status(http.StatusNotImplemented)
	return

	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Failed to parse multipart form",
			"requestID": requestID,
		})
		return
	}

	files := form.File["files"]
	if len(files) <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Need at least 2 files",
			"requestID": requestID,
		})
		return
	}

	maxUploads, _ := strconv.Atoi(os.Getenv("UPLOAD_BULK_MAX"))

	if len(files) > maxUploads {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Too many files provided",
			"requestID": requestID,
		})
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	var sum atomic.Int64

	// A single error file kills the entire request
	error := make(chan error, 1)

	for idx, header := range files {
		go func(fh *multipart.FileHeader, idx int) {
			defer wg.Done()

			select {
			case <-error:
				zap.L().Debug("Goroutine killed due to errors", zap.Int("goroutine_#", idx))
				return
			case <-c.Request.Context().Done():
				return
			default:
			}

			size, err := fileValidator(fh)
			if err != nil {
				error <- err
				return
			}

			sum.Add(size)
		}(header, idx)
	}

	wg.Wait()

	if err := <-error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"requestID": requestID,
		})
		return
	}

	var maxStorage, usedStorage int64

	err = d.DB.Gorm.
		Model(model.Stats{}).
		Where("user_id = ?", userID).
		Pluck("max_storage", &maxStorage).
		Pluck("used_storage", &usedStorage).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to pluck limits from db", zap.String("requestID", requestID), zap.String("userID", userID), zap.Error(err))
		return
	}

	if usedStorage+sum.Load() > maxStorage {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error":     "Insufficient storage",
			"requestID": requestID,
		})
		return
	}

	// for _, fh := range files {
	// 	d.Uploader.Do()
	// }
}

func fileValidator(fh *multipart.FileHeader) (int64, error) {
	maxFileSize, _ := strconv.ParseInt(os.Getenv("UPLOAD_MAX_SIZE"), 10, 64)

	if fh.Size > maxFileSize {
		return 0, fmt.Errorf("file %q is too large, max size is %d", fh.Filename, maxFileSize)
	}

	if !strings.HasPrefix(fh.Header.Get("Content-Type"), "video/") {
		return 0, fmt.Errorf("file %q is not a valid video type", fh.Filename)
	}

	f, err := fh.Open()
	if err != nil {
		return 0, fmt.Errorf("failed to open file %q: %v", fh.Filename, err)
	}
	defer f.Close()

	lim := io.LimitReader(f, maxFileSize)
	n, err := io.Copy(io.Discard, lim)
	if err != nil {
		return 0, fmt.Errorf("failed to read file %q: %v", fh.Filename, err)
	}

	if n >= maxFileSize {
		return 0, fmt.Errorf("file %q is too large, max size is %d", fh.Filename, maxFileSize)
	}

	mime, err := mimetype.DetectReader(f)
	if err != nil {
		return 0, fmt.Errorf("failed to detect mime type for file %q: %v", fh.Filename, err)
	}

	if !strings.HasPrefix(mime.String(), "video/") {
		return 0, fmt.Errorf("file %q is not a valid video type", fh.Filename)
	}

	return n, nil
}
