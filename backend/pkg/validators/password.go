package validators

import "errors"

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	ErrPasswordInvalid  = errors.New("password contains invalid characters")
	ErrPasswordTooLong  = errors.New("password is too long")
	ErrPasswordEmpty    = errors.New("no password provided")
)

func PasswordValidator(p string) error {
	if p == "" {
		return ErrPasswordEmpty
	}

	if len(p) < 8 {
		return ErrPasswordTooShort
	}

	if len(p) > 255 {
		return ErrPasswordTooLong
	}

	return nil
}
