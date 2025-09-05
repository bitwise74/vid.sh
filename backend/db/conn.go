// Package db contains things related to SQlite
package db

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/pkg/util"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() (*gorm.DB, error) {
	zap.L().Debug("db.New() entry")
	// If running in a docker container don't allow the sqlite file to be created.
	// The host should instead mount it using volumes
	if util.IsRunningInDocker() {
		zap.L().Debug("Running in Docker, checking for mounted database.db")
		if _, err := os.Stat("database.db"); err != nil {
			if err == os.ErrNotExist {
				zap.L().Error("SQLite database file not mounted")
				return nil, fmt.Errorf("SQLite database file not mounted, please use docker volumes to mount it to /app/database.db")
			}
		}
	}

	zap.L().Debug("Opening SQLite database")
	db, err := gorm.Open(sqlite.Open("database.db"))
	if err != nil {
		zap.L().Error("Failed to initialize SQLite database", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize SQLite database, %w", err)
	}

	zap.L().Debug("Running automigrate on tables")
	err = db.AutoMigrate(model.User{}, model.File{}, model.Stats{}, model.VerificationToken{}, model.ResendRequest{}, model.Migration{})
	if err != nil {
		zap.L().Error("Failed to automigrate tables", zap.Error(err))
		return nil, fmt.Errorf("failed to automigrate tables, %w", err)
	}

	zap.L().Debug("db.New() exit")
	return db, nil
}
