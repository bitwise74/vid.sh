package migrations

import (
	"bitwise74/video-api/internal/model"

	"gorm.io/gorm"
)

type CleanupTables struct{}

func (c CleanupTables) Name() string {
	return "001-cleanup-tables"
}

func (c CleanupTables) Exec(d *gorm.DB) error {
	err := DelCol(model.User{}, "verification_tokens", d)
	if err != nil {
		return err
	}

	err = DelCol(model.User{}, "resend_requests", d)
	if err != nil {
		return err
	}

	err = DelCol(model.Stats{}, "total_views", d)
	if err != nil {
		return err
	}

	err = DelCol(model.Stats{}, "total_watchtime", d)
	if err != nil {
		return err
	}

	err = DelTab(model.ResendRequest{}, d)
	if err != nil {
		return err
	}

	err = DelCol(model.File{}, "thumb_key", d)
	if err != nil {
		return err
	}

	err = DelCol(model.File{}, "views", d)
	if err != nil {
		return err
	}

	return nil
}
