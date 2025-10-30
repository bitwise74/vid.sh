package file

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/types"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type filterOpts struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
	Tags  string `json:"tags"`
}

var (
	validSortOpts   = []string{"newest", "oldest", "az", "za", "size-asc", "size-desc"}
	validResultOpts = []int{10, 20, 50, 100, 250} // Same as limit options
)

func FetchBulk(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	data := filterOpts{
		Limit: 10,
		Page:  1, // Because we already load one page initially
		Sort:  "newest",
		Tags:  "",
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Malformed or invalid query string",
			"requestID": requestID,
		})

		zap.L().Warn("Failed to bind fetch bulk query", zap.Error(err))
		return
	}

	if !slices.Contains(validSortOpts, data.Sort) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid sorting option",
			"requestID": requestID,
		})
		return
	}

	if !slices.Contains(validResultOpts, data.Limit) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid limit option",
			"requestID": requestID,
		})
		return
	}

	page := max(data.Page, 1)
	limit := min(max(data.Limit, 1), 250)
	sort := data.Sort

	fmt.Println(page, limit, sort)

	order := "created_at desc"
	switch sort {
	case "newest":
		order = "created_at desc"
	case "oldest":
		order = "created_at asc"
	case "az":
		order = "original_filename asc"
	case "za":
		order = "original_filename desc"
	case "size-asc":
		order = "file_size asc"
	case "size-desc":
		order = "file_size desc"
	}

	var entries []model.File

	err := d.DB.Gorm.
		Where("user_id = ?", userID).
		Order(order).
		Limit(limit).
		Offset(limit * page).
		Find(&entries).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"requestID": requestID,
		})
		zap.L().Error("Failed to lookup user files", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, entries)
}
