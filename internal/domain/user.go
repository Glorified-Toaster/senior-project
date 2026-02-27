package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines the possible roles for a user in the system.
type UserRole string

const (
	RoleStudent    UserRole = "STUDENT"
	RoleInstructor UserRole = "INSTRUCTOR"
	RoleAdmin      UserRole = "ADMIN"
)

// User represents a user entity in the examination system domain.
type User struct {
	ID           uuid.UUID
	Username     string
	FullName     string
	PasswordHash string
	Role         UserRole
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// IsValidRole checks if the provided role is valid within the system.
func IsValidRole(role UserRole) bool {
	switch role {
	case RoleStudent, RoleInstructor, RoleAdmin:
		return true
	default:
		return false
	}
}
