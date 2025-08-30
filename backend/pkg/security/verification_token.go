package security

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/pkg/util"
	"errors"
	"time"
)

const (
	tokenSize = 32
)

type VerificationTokenOpts struct {
	UserID    string
	Purpose   string
	ExpiresAt *time.Time
	CleanupAt *time.Time
}

func MakeVerificationToken(o *VerificationTokenOpts) (*model.VerificationToken, error) {
	if o == nil {
		return nil, errors.New("no token options provided")
	}

	if o.UserID == "" {
		return nil, errors.New("no user ID provided")
	}

	if o.Purpose == "" {
		return nil, errors.New("no token purpose provided")
	}

	if o.ExpiresAt == nil {
		return nil, errors.New("no expiry provided")
	}

	token, err := util.GenerateToken(tokenSize)
	if err != nil {
		return nil, err
	}

	return &model.VerificationToken{
		UserID:    o.UserID,
		Token:     token,
		Purpose:   o.Purpose,
		ExpiresAt: *o.ExpiresAt,
		CreatedAt: time.Now(),
		CleanupAt: o.CleanupAt,
		Used:      false,
	}, nil
}
