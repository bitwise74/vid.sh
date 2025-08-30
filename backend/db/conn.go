// Package db contains things related to SQlite
package db

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/pkg/util"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() (*gorm.DB, error) {
	// If running in a docker container don't allow the sqlite file to be created.
	// The host should instead mount it using volumes
	if util.IsRunningInDocker() {
		if _, err := os.Stat("database.db"); err != nil {
			if err == os.ErrNotExist {
				return nil, fmt.Errorf("SQLite database file not mounted, please use docker volumes to mount it to /app/database.db")
			}
		}
	}

	db, err := gorm.Open(sqlite.Open("database.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite database, %w", err)
	}

	err = db.AutoMigrate(model.User{}, model.File{}, model.Stats{}, model.VerificationToken{}, model.ResendRequest{}, model.Migration{})
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate tables, %w", err)
	}

	return db, nil
}
