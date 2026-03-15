package auth

import "errors"

var (
	ErrMissingToken error = errors.New("missing token")
	ErrInvalidToken error = errors.New("invalid token")
)
