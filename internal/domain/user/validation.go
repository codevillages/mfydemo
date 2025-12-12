package user

import (
	"errors"
	"net/mail"
	"unicode/utf8"
)

// ValidateNew validates fields for creating a user.
func ValidateNew(username, password, email string, status int32) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validateStatus(status); err != nil {
		return err
	}
	return nil
}

// ValidateUpdate validates fields for updating a user.
func ValidateUpdate(password, email string, status int32) error {
	if password != "" {
		if err := validatePassword(password); err != nil {
			return err
		}
	}
	if email != "" {
		if err := validateEmail(email); err != nil {
			return err
		}
	}
	if err := validateStatus(status); err != nil {
		return err
	}
	return nil
}

func validateUsername(username string) error {
	length := utf8.RuneCountInString(username)
	if length < 3 || length > 64 {
		return ErrInvalidAccount
	}
	return nil
}

func validatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < 8 || length > 128 {
		return ErrInvalidAccount
	}
	return nil
}

func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}

func validateStatus(status int32) error {
	switch status {
	case StatusInactive, StatusActive, StatusBlocked:
		return nil
	default:
		return ErrInvalidStatus
	}
}

// HashPassword abstracts hashing so infra can inject implementation.
type HashPassword func(raw string) (string, error)

// CompareHash abstracts password compare.
type CompareHash func(raw, hashed string) error

var (
	ErrHashPassword = errors.New("hash password error")
)
