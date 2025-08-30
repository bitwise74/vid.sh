package file

import (
	"bitwise74/video-api/internal"
	"bitwise74/video-api/internal/model"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var validLimits = []int{10, 20, 50, 100, 250}

func FileSearch(c *gin.Context, d *internal.Deps) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	searchQuery := strings.ToLower(c.Query("query"))
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No search query provided",
			"requestID": requestID,
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || !slices.Contains(validLimits, limit) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid limit provided",
			"requestID": requestID,
		})
		return
	}

	pageStr := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid page provided",
			"requestID": requestID,
		})
		return
	}

	var results []model.File

	err = d.DB.
		Where("user_id = ? AND original_name LIKE ?", userID, "%"+searchQuery+"%").
		Order("created_at desc").
		Offset(page * limit).
		Limit(limit).
		Find(&results).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})

		zap.L().Error("Failed to find files by search query", zap.Error(err))
		return
	}

	for i, file := range results {
		version := strconv.Itoa(file.Version)
		results[i].FileKey = file.FileKey + "?v=" + version
		results[i].ThumbKey = file.ThumbKey + "?v=" + version
	}

	c.JSON(http.StatusOK, results)
}
