package user

import "errors"

var (
	ErrUserExists     = errors.New("user already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidStatus  = errors.New("invalid status")
	ErrInvalidEmail   = errors.New("invalid email")
	ErrInvalidAccount = errors.New("invalid username or password")
)
