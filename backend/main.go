package main

import (
	"bitwise74/video-api/app"
	"bitwise74/video-api/config"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	err := config.Setup()
	if err != nil {
		panic(err)
	}

	router, err := app.NewRouter()
	if err != nil {
		panic(err)
	}

	zap.L().Info("Server starting")

	err = router.Run(":" + os.Getenv("HOST_PORT"))
	if err != nil {
		panic(err)
	}
}
