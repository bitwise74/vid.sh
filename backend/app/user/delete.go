package user

import (
	"bitwise74/video-api/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context, d *types.Dependencies) {
	c.Status(http.StatusNotImplemented)
	// requestID := c.MustGet("requestID").(string)
	// userID := c.MustGet("userID").(string)

	// var user model.User

	// err := d.DB.
	// 	Where("id = ?", userID).
	// 	Error
	// if err != nil {

	// }

	// resp, err := d.S3.C.GetObject()
}
