// Package repository provides implementations of the domain's repository interfaces.
package repository

import (
	"context"
	"errors"

	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/database/sqlc"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
	"uot-exam/internal/utils/password"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queries *sqlc.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

// Create : create a new user
func (r *UserRepository) Create(ctx context.Context, arg ports.CreateUserParams) (domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	if !domain.IsValidRole(arg.Role) {
		return domain.User{}, domain.ErrInvalidRole
	}

	if err := password.ValidatePasswordCriteria(arg.Password); err != nil {
		return domain.User{}, err
	}

	passwdHash, err := password.Hash(arg.Password)
	if err != nil {
		return domain.User{}, err
	}

	user, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username:     arg.Username,
		FullName:     arg.FullName,
		PasswordHash: string(passwdHash),
		Role:         sqlc.UserRoleType(arg.Role),
		IsActive:     arg.IsActive,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.User{}, domain.ErrUserAlreadyExists
		}
		return domain.User{}, err
	}

	return domain.User{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Role:      domain.UserRole(user.Role),
		IsActive:  user.IsActive,
		LastLogin: toTimePtr(user.LastLogin),
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		DeletedAt: toTimePtr(user.DeletedAt),
	}, nil
}

// GetByUsername : get a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	user, err := queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return domain.User{
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		PasswordHash: user.PasswordHash,
		Role:         domain.UserRole(user.Role),
		IsActive:     user.IsActive,
		LastLogin:    toTimePtr(user.LastLogin),
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    toTimePtr(user.DeletedAt),
	}, nil
}

// Login : login a user
func (r *UserRepository) Login(ctx context.Context, arg ports.LoginParams) (domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	user, err := queries.GetUserByUsername(ctx, arg.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	if !user.IsActive {
		return domain.User{}, domain.ErrUserNotActive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(arg.Password)); err != nil {
		return domain.User{}, domain.ErrInvalidPassword
	}

	if err := queries.UpdateUserLastLogin(ctx, user.ID); err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Role:      domain.UserRole(user.Role),
		IsActive:  user.IsActive,
		LastLogin: toTimePtr(user.LastLogin),
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		DeletedAt: toTimePtr(user.DeletedAt),
	}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	user, err := queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return domain.User{
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		PasswordHash: user.PasswordHash,
		Role:         domain.UserRole(user.Role),
		IsActive:     user.IsActive,
		LastLogin:    toTimePtr(user.LastLogin),
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    toTimePtr(user.DeletedAt),
	}, nil
}

func (r *UserRepository) UpdateUserInfo(ctx context.Context, arg ports.UpdateUserInfoParams) (domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	user, err := queries.UpdateUserInfo(ctx, sqlc.UpdateUserInfoParams{
		ID:       arg.ID,
		Username: arg.Username,
		FullName: arg.FullName,
		Role:     sqlc.UserRoleType(arg.Role),
		IsActive: arg.IsActive,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.User{}, domain.ErrUserAlreadyExists
		}
		return domain.User{}, err
	}

	return domain.User{
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		PasswordHash: user.PasswordHash,
		Role:         domain.UserRole(user.Role),
		IsActive:     user.IsActive,
		LastLogin:    toTimePtr(user.LastLogin),
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    toTimePtr(user.DeletedAt),
	}, nil
}

// SoftDelete : soft delete a user
func (r *UserRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	err := queries.SoftDeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Delete : delete a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	err := queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Restore : restore a user
func (r *UserRepository) Restore(ctx context.Context, id uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	err := queries.RestoreUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Disable : disable a user
func (r *UserRepository) Disable(ctx context.Context, id uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	err := queries.DisableUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Enable : enable a user
func (r *UserRepository) Enable(ctx context.Context, id uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	err := queries.EnableUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// ListDeleted : list all deleted users
func (r *UserRepository) ListDeleted(ctx context.Context, arg ports.ListDeletedUsersParams) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	users, err := queries.ListDeletedUsers(ctx, sqlc.ListDeletedUsersParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      domain.UserRole(user.Role),
			IsActive:  user.IsActive,
			LastLogin: toTimePtr(user.LastLogin),
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			DeletedAt: toTimePtr(user.DeletedAt),
		})
	}

	return domainUsers, nil
}

// ListAll : list all users
func (r *UserRepository) ListAll(ctx context.Context, arg ports.ListAllUsersParams) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	users, err := queries.ListAllUsers(ctx, sqlc.ListAllUsersParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      domain.UserRole(user.Role),
			IsActive:  user.IsActive,
			LastLogin: toTimePtr(user.LastLogin),
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			DeletedAt: toTimePtr(user.DeletedAt),
		})
	}

	return domainUsers, nil
}

func (r *UserRepository) Search(ctx context.Context, arg ports.SearchUsersParams) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	users, err := queries.SearchUsers(ctx, sqlc.SearchUsersParams{
		Search: arg.Search,
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})

	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      domain.UserRole(user.Role),
			IsActive:  user.IsActive,
			LastLogin: toTimePtr(user.LastLogin),
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			DeletedAt: toTimePtr(user.DeletedAt),
		})
	}

	return domainUsers, nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	count, err := queries.CountUsers(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) CountDeleted(ctx context.Context) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	count, err := queries.CountDeletedUsers(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) CountSearch(ctx context.Context, search string) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	count, err := queries.CountSearchUsers(ctx, search)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) CountSearchDeleted(ctx context.Context, search string) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	count, err := queries.CountSearchDeletedUsers(ctx, search)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) SearchDeleted(ctx context.Context, arg ports.SearchUsersParams) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	users, err := queries.SearchDeletedUsers(ctx, sqlc.SearchDeletedUsersParams{
		Search: arg.Search,
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      domain.UserRole(user.Role),
			IsActive:  user.IsActive,
			LastLogin: toTimePtr(user.LastLogin),
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			DeletedAt: toTimePtr(user.DeletedAt),
		})
	}

	return domainUsers, nil
}

func (r *UserRepository) ListAllInstructors(ctx context.Context, arg ports.ListAllInstructorsParams) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	users, err := queries.ListAllInstructors(ctx, sqlc.ListAllInstructorsParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      domain.UserRole(user.Role),
			IsActive:  user.IsActive,
			LastLogin: toTimePtr(user.LastLogin),
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			DeletedAt: toTimePtr(user.DeletedAt),
		})
	}

	return domainUsers, nil
}

func (r *UserRepository) ListAllStudents(ctx context.Context, arg ports.ListAllStudentsParams) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	users, err := queries.ListAllStudents(ctx, sqlc.ListAllStudentsParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      domain.UserRole(user.Role),
			IsActive:  user.IsActive,
			LastLogin: toTimePtr(user.LastLogin),
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			DeletedAt: toTimePtr(user.DeletedAt),
		})
	}

	return domainUsers, nil
}

func (r *UserRepository) CountInstructors(ctx context.Context) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	count, err := queries.CountInstructors(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) CountStudents(ctx context.Context) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	count, err := queries.CountStudents(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) SearchStudents(ctx context.Context, arg ports.SearchStudentsParams) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	users, err := queries.SearchStudents(ctx, sqlc.SearchStudentsParams{
		Search: arg.Search,
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      domain.UserRole(user.Role),
			IsActive:  user.IsActive,
			LastLogin: toTimePtr(user.LastLogin),
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			DeletedAt: toTimePtr(user.DeletedAt),
		})
	}

	return domainUsers, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: hashedPassword,
	})
}

func (r *UserRepository) UpdateLastInteraction(ctx context.Context, id uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.UpdateUserLastLogin(ctx, id)
}
