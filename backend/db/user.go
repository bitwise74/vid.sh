package db

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/pkg/security"
	"os"
	"strconv"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	validationDeadline = 24 * 7 * time.Hour // 1 Week
)

// CreateUserWithToken creates a new user and a verification token for that user.
// The user is created as unverified and the token is set to expire in 1 week.
// The token is also set to be cleaned up in 30 days.
func (d *DB) CreateUserWithToken(email, passwd string) (*model.User, error) {
	userID := gonanoid.MustGenerate(charset, 10)

	passwdHash, err := security.Argon.GenerateFromPassword(passwd)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(validationDeadline)
	cleanupAt := time.Now().Add(time.Hour * 24 * 30) // Month
	maxStorage, _ := strconv.ParseInt(os.Getenv("STORAGE_MAX_USAGE"), 10, 64)

	token, err := security.MakeToken(&security.TokenOpts{
		UserID:    userID,
		Purpose:   "verification",
		ExpiresAt: &expiresAt,
		CleanupAt: &cleanupAt,
	})
	if err != nil {
		return nil, err
	}

	u := &model.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passwdHash,
		Verified:     false,
		ExpiresAt:    &expiresAt,
		CreatedAt:    time.Now(),
		Tokens:       []model.Token{*token},
		Stats: model.Stats{
			UserID:     userID,
			MaxStorage: maxStorage,
		},
	}

	err = d.Gorm.Create(u).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return nil, ErrUserExists
		}

		return nil, err
	}

	return u, nil
}

// FetchUserDataByID fetches a user and their basic information by their ID.
func (d *DB) FetchUserDataByID(id string) (*model.User, error) {
	var user model.User

	err := d.Gorm.
		Where("id = ?", id).
		Select("id", "username", "avatar_hash", "public_profile_enabled").
		Preload("Files", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10).Order("created_at DESC")
		}).
		Preload("Stats").
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FetchUserByEmail fetches a user by their email.
func (d *DB) FetchUserByEmail(email string) (*model.User, error) {
	var user model.User

	err := d.Gorm.
		Where("email = ?", email).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// PluckUserID fetches and returns a user's ID by their email.
func (d *DB) PluckUserID(email string) (string, error) {
	var userID string

	err := d.Gorm.
		Model(&model.User{}).
		Where("email = ?", email).
		Pluck("id", &userID).
		Error
	if err != nil {
		return "", err
	}

	return userID, nil
}

// CheckIfUsernameTaken checks if a username is already taken.
func (d *DB) CheckIfUsernameTaken(username string) (bool, error) {
	// This will be checked further later so its fine to return false
	if username == "" {
		return false, nil
	}

	var count int64

	err := d.Gorm.
		Model(&model.User{}).
		Where("username = ?", username).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
