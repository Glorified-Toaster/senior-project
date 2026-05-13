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
	UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (domain.User, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
	Enable(ctx context.Context, id uuid.UUID) error
	ListDeleted(ctx context.Context, arg ListDeletedUsersParams) ([]domain.User, error)
	ListAll(ctx context.Context, arg ListAllUsersParams) ([]domain.User, error)
	Search(ctx context.Context, arg SearchUsersParams) ([]domain.User, error)
	SearchDeleted(ctx context.Context, arg SearchUsersParams) ([]domain.User, error)
	Count(ctx context.Context) (int64, error)
	CountDeleted(ctx context.Context) (int64, error)
	CountSearch(ctx context.Context, search string) (int64, error)
	CountSearchDeleted(ctx context.Context, search string) (int64, error)
	ListAllInstructors(ctx context.Context, arg ListAllInstructorsParams) ([]domain.User, error)
	ListAllStudents(ctx context.Context, arg ListAllStudentsParams) ([]domain.User, error)
	CountInstructors(ctx context.Context) (int64, error)
	CountStudents(ctx context.Context) (int64, error)
	SearchStudents(ctx context.Context, arg SearchStudentsParams) ([]domain.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error
	UpdateLastInteraction(ctx context.Context, id uuid.UUID) error
}

type CreateUserParams struct {
	Username string          `json:"username"  form:"username"  validate:"required,min=3,max=20"`
	FullName string          `json:"full_name" form:"full_name" validate:"required,min=3,max=20"`
	Password string          `json:"password"  form:"password"  validate:"required"`
	Role     domain.UserRole `json:"role"      form:"role"      validate:"required"`
	IsActive bool            `json:"is_active" form:"is_active"`
}

type UpdateUserInfoParams struct {
	ID       uuid.UUID       `json:"id"        form:"id"        validate:"required"`
	Username string          `json:"username"  form:"username"  validate:"required,min=3,max=20"`
	FullName string          `json:"full_name" form:"full_name" validate:"required,min=3,max=20"`
	Role     domain.UserRole `json:"role"      form:"role"      validate:"required"`
	IsActive bool            `json:"is_active" form:"is_active"`
}

type LoginParams struct {
	Username string `json:"username" form:"username" validate:"required,min=4,max=20"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=128"`
}

type ListAllUsersParams struct {
	Limit  int32 `json:"limit" validate:"required"`
	Offset int32 `json:"offset" validate:"required"`
}

type SearchUsersParams struct {
	Search string `json:"search" validate:"required"`
	Limit  int32  `json:"limit" validate:"required"`
	Offset int32  `json:"offset" validate:"required"`
}

type ListDeletedUsersParams struct {
	Limit  int32 `json:"limit" validate:"required"`
	Offset int32 `json:"offset" validate:"required"`
}

type ListAllStudentsParams struct {
	Limit  int32 `json:"limit" validate:"required"`
	Offset int32 `json:"offset" validate:"required"`
}

type SearchStudentsParams struct {
	Search string `json:"search" validate:"required"`
	Limit  int32  `json:"limit" validate:"required"`
	Offset int32  `json:"offset" validate:"required"`
}
