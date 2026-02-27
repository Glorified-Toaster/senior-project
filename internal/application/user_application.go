package application

import (
	"context"

	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
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
