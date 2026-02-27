package domain

import "errors"

var (
	// User errors (postgres)
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	// JWT errors
	ErrInvalidJWTToken = errors.New("invalid token")
	ErrExpiredJWTToken = errors.New("token has expired")
	ErrMissingJWTKey   = errors.New("JWT key not set")
)
