package ports

import (
	"context"
	"uot-exam/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, arg CreateUserParams) (domain.User, error)
	Login(ctx context.Context, arg LoginParams) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
	Enable(ctx context.Context, id uuid.UUID) error
	ListDeleted(ctx context.Context) ([]domain.User, error)
	ListAll(ctx context.Context) ([]domain.User, error)
}

type CreateUserParams struct {
	Username string          `json:"username" validate:"required,min=3,max=20"`
	FullName string          `json:"full_name" validate:"required,min=3,max=20"`
	Password string          `json:"password" validate:"required"`
	Role     domain.UserRole `json:"role" validate:"required"`
	IsActive bool            `json:"is_active" validate:"required"`
}

type LoginParams struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=3,max=128"`
}
