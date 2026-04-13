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

func (app *Application) CountSearchExams(ctx context.Context, search string) (int64, error) {
	var count int64
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		count, err = app.examRepo.CountSearch(txCtx, search)
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

func (app *Application) UpdateExam(ctx context.Context, arg ports.UpdateExamParams) (domain.Exam, error) {
	var exam domain.Exam
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exam, err = app.examRepo.Update(txCtx, arg)
		return err
	})
	return exam, err
}

func (app *Application) SearchExamsBySubject(ctx context.Context, arg ports.SearchExamsBySubjectParams) ([]domain.Exam, error) {
	var exams []domain.Exam
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exams, err = app.examRepo.SearchBySubject(txCtx, arg)
		return err
	})
	return exams, err
}

func (app *Application) CountSearchExamsBySubject(ctx context.Context, subjectID uuid.UUID, search string) (int64, error) {
	var count int64
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		count, err = app.examRepo.CountSearchBySubject(txCtx, subjectID, search)
		return err
	})
	return count, err
}

func (app *Application) ListExamsForStudent(ctx context.Context, studentID uuid.UUID) ([]domain.Exam, error) {
	var exams []domain.Exam
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exams, err = app.examRepo.ListExamsForStudent(txCtx, studentID)
		return err
	})
	return exams, err
}

func (app *Application) ListExamsCreatedBy(ctx context.Context, instructorID uuid.UUID) ([]domain.Exam, error) {
	var exams []domain.Exam
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		exams, err = app.examRepo.ListExamsCreatedBy(txCtx, instructorID)
		return err
	})
	return exams, err
}

func (app *Application) ListAttemptsByStudent(ctx context.Context, studentID uuid.UUID) ([]domain.ExamAttempt, error) {
	var attempts []domain.ExamAttempt
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		attempts, err = app.examRepo.ListAttemptsByStudent(txCtx, studentID)
		return err
	})
	return attempts, err
}

