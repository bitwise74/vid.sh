package security

import (
	"bitwise74/video-api/internal/model"
	"bitwise74/video-api/pkg/util"
	"fmt"
	"time"
)

const (
	tokenSize = 32
)

type TokenOpts struct {
	UserID    string
	Purpose   string
	ExpiresAt *time.Time
	CleanupAt *time.Time
}

// MakeToken creates a new token for a user with the specified options.
// WARNING: This function has no sanity checks, you can put in ridiculous values and it will
// create a token with those values. Make sure to validate the options before calling this function.
// All this does is abstract away the token generation and struct creation.
func MakeToken(o *TokenOpts) (*model.Token, error) {
	token, err := util.GenerateToken(tokenSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token, %w", err)
	}

	return &model.Token{
		UserID:    o.UserID,
		Token:     token,
		Purpose:   o.Purpose,
		ExpiresAt: *o.ExpiresAt,
		CreatedAt: time.Now(),
		CleanupAt: o.CleanupAt,
		Used:      false,
	}, nil
}
