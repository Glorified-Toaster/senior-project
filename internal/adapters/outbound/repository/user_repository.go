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
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		PasswordHash: user.PasswordHash,
		Role:         domain.UserRole(user.Role),
		IsActive:     user.IsActive,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    &user.DeletedAt.Time,
	}, nil
}

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
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    &user.DeletedAt.Time,
	}, nil
}

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
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		PasswordHash: user.PasswordHash,
		Role:         domain.UserRole(user.Role),
		IsActive:     user.IsActive,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    &user.DeletedAt.Time,
		LastLogin:    &user.LastLogin.Time,
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
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    &user.DeletedAt.Time,
	}, nil
}

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

func (r *UserRepository) ListDeleted(ctx context.Context) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	users, err := queries.ListDeletedUsers(ctx)
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:           user.ID,
			Username:     user.Username,
			FullName:     user.FullName,
			PasswordHash: user.PasswordHash,
			Role:         domain.UserRole(user.Role),
			IsActive:     user.IsActive,
			CreatedAt:    user.CreatedAt.Time,
			UpdatedAt:    user.UpdatedAt.Time,
			DeletedAt:    &user.DeletedAt.Time,
		})
	}

	return domainUsers, nil
}

func (r *UserRepository) ListAll(ctx context.Context) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	users, err := queries.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:           user.ID,
			Username:     user.Username,
			FullName:     user.FullName,
			PasswordHash: user.PasswordHash,
			Role:         domain.UserRole(user.Role),
			IsActive:     user.IsActive,
			CreatedAt:    user.CreatedAt.Time,
			UpdatedAt:    user.UpdatedAt.Time,
			DeletedAt:    &user.DeletedAt.Time,
		})
	}

	return domainUsers, nil
}
