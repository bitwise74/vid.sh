// Package validators contains validators found throughout the application
// that have been abstracted away from the main code
package validators

import (
	"errors"
	"net/mail"
)

var (
	ErrEmailEmpty   = errors.New("no email address provided")
	ErrEmailInvalid = errors.New("invalid email address provided")
)

func EmailValidator(e string) error {
	if e == "" {
		return ErrEmailEmpty
	}

	if _, err := mail.ParseAddress(e); err != nil {
		return ErrEmailInvalid
	}

	return nil
}
