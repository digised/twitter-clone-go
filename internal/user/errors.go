package user

import "errors"

var (
	ErrNotFound      = errors.New("user not found")
	ErrUsernameTaken = errors.New("username already taken")
	ErrEmailTaken    = errors.New("email already taken")
)
