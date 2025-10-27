// Package db contains things related to SQlite
package db

import (
	"bitwise74/video-api/internal/migrations"
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/pkg/util"
	"fmt"
	"os"
	"slices"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	Gorm *gorm.DB
}

var m = []any{
	migrations.CleanupTables{},
}

func New() (*DB, error) {
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

	err = db.AutoMigrate(model.User{}, model.File{}, model.Stats{}, model.Migration{}, model.Token{})
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate tables, %w", err)
	}

	// Run migrations
	var ranMigrations []string

	err = db.Model(model.Migration{}).Where("1 = 1").Select("name").Find(&ranMigrations).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get last migration name, %w", err)
	}

	for _, migration := range m {
		name := migration.(migrations.CleanupTables).Name()
		if slices.Contains(ranMigrations, name) {
			continue
		}

		err = migration.(migrations.CleanupTables).Exec(db)
		if err != nil {
			return nil, fmt.Errorf("failed to run migration %s, %w", name, err)
		}

		err = db.Create(&model.Migration{
			Name: name,
		}).Error
		if err != nil {
			return nil, fmt.Errorf("failed to record migration %s, %w", name, err)
		}

		zap.L().Info("Ran migration", zap.String("name", name))
	}

	return &DB{
		Gorm: db,
	}, nil
}
