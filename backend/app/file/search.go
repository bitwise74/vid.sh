package file

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SearchRequestBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
}

var validLimits = []int{10, 20, 50, 100, 250}

func Search(c *gin.Context, d *types.Dependencies) {
	requestID := c.MustGet("requestID").(string)
	userID := c.MustGet("userID").(string)

	var data SearchRequestBody
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid JSON body",
			"requestID": requestID,
		})
		return
	}

	var results []model.File

	err := d.DB.Gorm.
		Where("user_id = ? AND original_name LIKE ?", userID, "%"+data.Query+"%").
		Order("created_at desc").
		Offset(data.Page * data.Limit).
		Limit(data.Limit).
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

	c.JSON(http.StatusOK, results)
}
