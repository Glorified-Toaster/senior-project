package application

import (
	"context"
	"mime/multipart"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"

	"github.com/google/uuid"
)

func (app *Application) CreateQuestion(ctx context.Context, arg ports.CreateQuestionParams) (domain.Question, error) {
	var question domain.Question
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		question, err = app.questionRepo.Create(txCtx, arg)
		return err
	})
	return question, err
}

func (app *Application) CreateChoice(ctx context.Context, arg ports.CreateChoiceParams) (domain.Choice, error) {
	var choice domain.Choice
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		choice, err = app.questionRepo.CreateChoice(txCtx, arg)
		return err
	})
	return choice, err
}

func (app *Application) ListQuestionsByExam(ctx context.Context, arg uuid.UUID) ([]domain.Question, error) {
	var questions []domain.Question
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		questions, err = app.questionRepo.ListQuestionsByExam(txCtx, arg)
		return err
	})
	return questions, err
}

func (app *Application) ListChoicesByQuestion(ctx context.Context, arg uuid.UUID) ([]domain.Choice, error) {
	var choices []domain.Choice
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		choices, err = app.questionRepo.ListChoicesByQuestion(txCtx, arg)
		return err
	})
	return choices, err
}

func (app *Application) GetQuestionByChecksum(ctx context.Context, arg string) (bool, error) {
	var exists bool
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exists, err = app.questionRepo.GetQuestionByChecksum(txCtx, arg)
		return err
	})
	return exists, err
}

func (app *Application) GetQuestionByID(ctx context.Context, arg uuid.UUID) (domain.Question, error) {
	var question domain.Question
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		question, err = app.questionRepo.GetByID(txCtx, arg)
		return err
	})
	return question, err
}

func (app *Application) DeleteChoicesByQuestion(ctx context.Context, arg uuid.UUID) error {
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		err = app.questionRepo.DeleteChoicesByQuestion(txCtx, arg)
		return err
	})
	return err
}

func (app *Application) DeleteQuestionAndChoices(ctx context.Context, arg uuid.UUID) error {
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		err = app.questionRepo.DeleteQuestionAndChoices(txCtx, arg)
		return err
	})
	return err
}

func (app *Application) UpdateQuestion(ctx context.Context, arg ports.UpdateQuestionParams) (domain.Question, error) {
	var question domain.Question
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		question, err = app.questionRepo.Update(txCtx, arg)
		return err
	})
	return question, err
}

func (app *Application) UpdateQuestionWithChoices(ctx context.Context, qArg ports.UpdateQuestionParams, cArgs []ports.CreateChoiceParams) (domain.Question, error) {
	var question domain.Question
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		question, err = app.questionRepo.Update(txCtx, qArg)
		if err != nil {
			return err
		}

		if len(cArgs) > 0 {
			// Fetch all existing choices (including soft-deleted)
			existingChoices, err := app.questionRepo.ListAllChoicesByQuestion(txCtx, qArg.ID)
			if err != nil {
				return err
			}

			// We use a fixed limit of 4 since the DB trigger enforces it
			for i := 0; i < 4; i++ {
				if i < len(cArgs) {
					// We want to have this choice as active
					if i < len(existingChoices) {
						// Update existing choice (this will also un-delete if it was soft-deleted)
						_, err = app.questionRepo.UpdateChoice(txCtx, ports.UpdateChoiceParams{
							ID:         existingChoices[i].ID,
							ChoiceText: cArgs[i].ChoiceText,
							IsCorrect:  cArgs[i].IsCorrect,
						})
					} else {
						// Create new choice
						_, err = app.questionRepo.CreateChoice(txCtx, cArgs[i])
					}
				} else {
					// We don't want this choice (extra ones)
					if i < len(existingChoices) {
						// Soft-delete it if it's not already deleted
						if existingChoices[i].DeletedAt == nil || existingChoices[i].DeletedAt.IsZero() {
							err = app.questionRepo.SoftDeleteChoiceByID(txCtx, existingChoices[i].ID)
						}
					}
				}
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return question, err
}

func (app *Application) UploadQuestionImage(questionImage *multipart.FileHeader, hashedFilename string) (string, error) {
	file, err := questionImage.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	return app.localDisk.UploadFile(file, hashedFilename)
}
