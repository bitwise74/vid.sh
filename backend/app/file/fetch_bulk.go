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
	"gorm.io/gorm"
)

// AZ = A - Z as in alphabetic same for ZA
var validSortOpts = []string{"newest", "oldest", "az", "za", "size-asc", "size-desc"}

func FileFetchBulk(c *gin.Context, d *internal.Deps) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	pageStr := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Page must be a number",
			"requestID": requestID,
		})
		return
	}

	if page < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Page can't be negative",
			"requestID": requestID,
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Limit must be a number",
			"requestID": requestID,
		})
		return
	}

	if limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Limit must be greater than 0",
			"requestID": requestID,
		})
		return
	}

	if limit > 250 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Limit must be smaller than 250",
			"requestID": requestID,
		})
		return
	}

	sort := strings.ToLower(c.DefaultQuery("sort", "newest"))
	if !slices.Contains(validSortOpts, sort) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid sorting option",
			"requestID": requestID,
		})
		return
	}

	order := ""

	switch sort {
	case "newest":
		order = "created_at desc"
	case "oldest":
		order = "created_at asc"
	case "az":
		order = "name"
	case "za":
		order = "name desc"
	case "size-asc":
		order = "size asc"
	case "size-desc":
		order = "size desc"
	}

	offset := page * limit
	var entries []model.File

	err = d.DB.
		Where("user_id = ?", userID).
		Order(order).
		Offset(offset).
		Limit(limit).
		Find(&entries).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "No results",
				"requestID": requestID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		zap.L().Error("Failed to lookup user files", zap.Error(err))
		return
	}

	for i, file := range entries {
		version := strconv.Itoa(file.Version)
		entries[i].FileKey = file.FileKey + "?v=" + version
		entries[i].ThumbKey = file.ThumbKey + "?v=" + version
	}

	c.JSON(http.StatusOK, entries)
}
