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

func (app *Application) Login(ctx context.Context, arg ports.LoginParams) (domain.User, error) {
	var user domain.User
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		user, err = app.userRepo.Login(txCtx, arg)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
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

func (app *Application) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		err = app.userRepo.SoftDelete(txCtx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) WipeUser(ctx context.Context, id uuid.UUID) error {
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		err = app.userRepo.Delete(txCtx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) RestoreUser(ctx context.Context, id uuid.UUID) error {
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		err = app.userRepo.Restore(txCtx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) DisableUser(ctx context.Context, id uuid.UUID) error {
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		err = app.userRepo.Disable(txCtx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) EnableUser(ctx context.Context, id uuid.UUID) error {
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		err = app.userRepo.Enable(txCtx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) ListDeletedUsers(ctx context.Context) ([]domain.User, error) {
	users, err := app.userRepo.ListDeleted(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (app *Application) ListAllUsers(ctx context.Context, arg ports.ListAllUsersParams) ([]domain.User, error) {
	users, err := app.userRepo.ListAll(ctx, arg)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (app *Application) SearchUsers(ctx context.Context, search string) ([]domain.User, error) {
	users, err := app.userRepo.Search(ctx, search)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (app *Application) CountUsers(ctx context.Context) (int64, error) {
	count, err := app.userRepo.Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}
