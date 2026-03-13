package users

import "errors"

var (
	ErrInvalidCredentials    error = errors.New("invalid credentials")
	ErrInvalidPlayerId       error = errors.New("invalid player ID")
	ErrPasswordTooShort      error = errors.New("password must be at least 8 characters long")
	ErrPlayerIdAlreadyExists error = errors.New("cannot use this username")
)
