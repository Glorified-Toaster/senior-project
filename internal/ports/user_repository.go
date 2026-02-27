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
	Username     string
	FullName     string
	PasswordHash string
	Role         domain.UserRole
	IsActive     bool
}
