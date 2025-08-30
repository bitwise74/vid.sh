package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type response struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

func NewTurnstileMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		turnstileEnabled, err := strconv.ParseBool(os.Getenv("TURNSTILE_ENABLED"))
		if err != nil {
			turnstileEnabled = false
		}

		if !turnstileEnabled {
			c.Next()
			return
		}

		token := c.Request.Header.Get("TurnstileToken")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing or invalid turnstile token",
			})
			return
		}

		payload := gin.H{
			"secret":   os.Getenv("TURNSTILE_SECRET_TOKEN"),
			"response": token,
			"remoteip": c.ClientIP(),
		}

		jsonBody, _ := json.Marshal(payload)
		resp, err := http.Post("https://challenges.cloudflare.com/turnstile/v0/siteverify", "application/json", bytes.NewReader(jsonBody))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		respBody, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		var res response
		if err := json.Unmarshal(respBody, &res); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		c.Next()
	}
}
