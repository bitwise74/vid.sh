package validators

import (
	"errors"
	"io"
	"mime/multipart"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrUsernameTooLong  = errors.New("username is too long")
	ErrUsernameTooShort = errors.New("username is too short")
	ErrUsernameInvalid  = errors.New("username is invalid")

	ErrAvatarFileTooBig = errors.New("avatar file is too big")
	ErrAvatarNotValid   = errors.New("provided avatar file is not an image")
	ErrAvatarEmpty      = errors.New("avatar file is empty")

	ErrNoAvatarCropOpts = errors.New("no avatar crop options provided")
	ErrNoUpdateOpts     = errors.New("no update options provided")

	ErrCropOptionsInvalid = errors.New("invalid cropping options provided")

	// Just some usernames that could be misleading/weird with javascript
	disallowedUsernames = []string{
		"admin",
		"root",
		"support",
		"system",
		"moderator",
		"staff",
		"help",
		"contact",
		"null",
		"undefined",
		"operator",
		"superuser",
		"vidsh",
	}

	allowedCharacters = `^[a-zA-Z0-9_-]+$`
)

type UserUpdateOpts struct {
	Username             string                `form:"username"`
	Avatar               *multipart.FileHeader `form:"avatar"`
	PublicProfileEnabled *bool                 `form:"publicProfileEnabled"`

	// In format X,Y,W,H
	AvatarCrop string `form:"avatarCrop"`

	AvatarCropProper []int

	DefaultPrivateVideos *bool `form:"defaultPrivateVideos"`
}

func (u UserUpdateOpts) IsEmpty() bool {
	return u.Username == "" && u.Avatar == nil && u.PublicProfileEnabled == nil && u.AvatarCrop == "" && u.DefaultPrivateVideos == nil
}

const maxAvatarSize = 5242880

func UserValidator(o *UserUpdateOpts) error {
	if o.IsEmpty() {
		return ErrNoUpdateOpts
	}

	if o.Username != "" {
		lower := strings.ToLower(o.Username)

		if len(o.Username) < 3 {
			return ErrUsernameTooShort
		}

		if len(o.Username) > 30 {
			return ErrUsernameTooLong
		}

		match, err := regexp.MatchString(allowedCharacters, o.Username)
		if err != nil || !match {
			return ErrUsernameInvalid
		}

		if slices.Contains(disallowedUsernames, lower) {
			return ErrUsernameInvalid
		}
	}

	if o.Avatar != nil {
		if o.Avatar.Size == 0 {
			return ErrAvatarEmpty
		}

		if o.Avatar.Size > maxAvatarSize {
			return ErrAvatarFileTooBig
		}

		if !strings.HasPrefix(o.Avatar.Header.Get("Content-Type"), "image/") {
			return ErrAvatarNotValid
		}

		f, err := o.Avatar.Open()
		if err != nil {
			return err
		}
		defer f.Close()

		// Check for size
		limited := io.LimitReader(f, maxAvatarSize)
		n, err := io.Copy(io.Discard, limited)
		if err != nil {
			return err
		}

		if n > maxAvatarSize {
			return ErrAvatarFileTooBig
		}
	}

	if o.AvatarCrop != "" {
		opts := strings.Split(o.AvatarCrop, ",")
		if len(opts) != 4 {
			return ErrCropOptionsInvalid
		}

		xStr, yStr, wStr, hStr := opts[0], opts[1], opts[2], opts[3]

		x, err := strconv.Atoi(xStr)
		if err != nil {
			return ErrCropOptionsInvalid
		}

		y, err := strconv.Atoi(yStr)
		if err != nil {
			return ErrCropOptionsInvalid
		}

		w, err := strconv.Atoi(wStr)
		if err != nil {
			return ErrCropOptionsInvalid
		}

		h, err := strconv.Atoi(hStr)
		if err != nil {
			return ErrCropOptionsInvalid
		}

		if x <= 0 || y <= 0 || w <= 0 || h <= 0 {
			return ErrCropOptionsInvalid
		}

		o.AvatarCropProper = append(o.AvatarCropProper, x, y, w, h)
	}

	return nil
}
