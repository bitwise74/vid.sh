package db

import "bitwise74/video-api/internal/model"

func (d *DB) CheckFileOwnership(userID, fileKey string) (bool, error) {
	var owns bool

	err := d.Gorm.
		Model(model.File{}).
		Where("user_id = ? AND file_key = ?", userID, fileKey).
		Select("count(*) > 0").
		Scan(&owns).
		Error
	if err != nil {
		return false, err
	}

	return owns, nil
}
