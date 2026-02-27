package application

import (
	"context"

	"uot-exam/internal/domain"
	"uot-exam/internal/ports"

	"github.com/google/uuid"
)

func (app *Application) CreateUser(ctx context.Context, arg ports.CreateUserParams) (domain.User, error) {
	var createdUser domain.User
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		createdUser, err = app.userRepo.Create(txCtx, arg)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.User{}, err
	}
	return createdUser, nil
}

func (app *Application) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := app.userRepo.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (app *Application) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	user, err := app.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
