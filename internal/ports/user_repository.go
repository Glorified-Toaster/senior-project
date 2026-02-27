package ports

import (
	"context"
	"uot-exam/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, arg CreateUserParams) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type CreateUserParams struct {
	Username     string          `json:"username" validate:"required,min=3,max=20"`
	FullName     string          `json:"full_name" validate:"required,min=3,max=20"`
	PasswordHash string          `json:"password_hash" validate:"required"`
	Role         domain.UserRole `json:"role" validate:"required"`
	IsActive     bool            `json:"is_active" validate:"required"`
}
