package main

import (
	"bitwise74/video-api/app"
	"bitwise74/video-api/config"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	zap.L().Debug("main() entry")
	gin.SetMode(gin.ReleaseMode)

	zap.L().Debug("Calling config.Setup()")
	err := config.Setup()
	if err != nil {
		zap.L().Fatal("Config setup failed", zap.Error(err))
		panic(err)
	}

	zap.L().Debug("Creating router")
	router, err := app.NewRouter()
	if err != nil {
		zap.L().Fatal("Router creation failed", zap.Error(err))
		panic(err)
	}

	zap.L().Info("Server starting", zap.String("port", os.Getenv("HOST_PORT")))

	err = router.Run(":" + os.Getenv("HOST_PORT"))
	if err != nil {
		zap.L().Fatal("Server failed to start", zap.Error(err))
		panic(err)
	}
	zap.L().Debug("main() exit")
}
