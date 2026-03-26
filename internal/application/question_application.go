package application

import (
	"context"
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
