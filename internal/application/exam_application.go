package application

import (
	"context"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"

	"github.com/google/uuid"
)

func (app *Application) ListAllExams(ctx context.Context, arg ports.ListAllExamsParams) ([]domain.Exam, error) {
	var exams []domain.Exam
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exams, err = app.examRepo.ListAll(txCtx, arg)
		return err
	})
	return exams, err
}

func (app *Application) SearchExams(ctx context.Context, arg ports.SearchExamsParams) ([]domain.Exam, error) {
	var exams []domain.Exam
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exams, err = app.examRepo.Search(txCtx, arg)
		return err
	})
	return exams, err
}

func (app *Application) CountExams(ctx context.Context) (int64, error) {
	var count int64
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		count, err = app.examRepo.Count(txCtx)
		return err
	})
	return count, err
}

func (app *Application) SoftDeleteExam(ctx context.Context, id uuid.UUID) error {
	return app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return app.examRepo.SoftDelete(txCtx, id)
	})
}

func (app *Application) GetExamByID(ctx context.Context, id uuid.UUID) (domain.Exam, error) {
	var exam domain.Exam
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exam, err = app.examRepo.GetByID(txCtx, id)
		return err
	})
	return exam, err
}

func (app *Application) CreateExam(ctx context.Context, arg ports.CreateExamParams) (domain.Exam, error) {
	var exam domain.Exam
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exam, err = app.examRepo.Create(txCtx, arg)
		return err
	})
	return exam, err
}

func (app *Application) ListExamsBySubject(ctx context.Context, subjectID uuid.UUID) ([]domain.Exam, error) {
	var exams []domain.Exam
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exams, err = app.examRepo.ListBySubject(txCtx, subjectID)
		return err
	})
	return exams, err
}
