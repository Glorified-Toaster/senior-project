package domain

import "errors"

var (
	// User errors (postgres)
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrAdminAlreadyExists = errors.New("an admin account already exists")
	ErrUserNotFound      = errors.New("user not found")
	// Exam errors
	ErrExamNotFound = errors.New("exam not found")

	// JWT errors
	ErrInvalidJWTToken = errors.New("invalid token")
	ErrExpiredJWTToken = errors.New("token has expired")
	ErrMissingJWTKey   = errors.New("JWT key not set")

	// Auth errors
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotActive   = errors.New("user is not active")

	// Role errors
	ErrInvalidRole = errors.New("invalid role")
)

type ErrCSVUpload struct {
	Message string
}

func (e *ErrCSVUpload) Error() string { return e.Message }
