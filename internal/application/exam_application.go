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

func (app *Application) GetAttemptByExamAndStudent(ctx context.Context, examID uuid.UUID, studentID uuid.UUID) (domain.ExamAttempt, error) {
	var attempt domain.ExamAttempt
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		attempt, err = app.examRepo.GetAttemptByExamAndStudent(txCtx, examID, studentID)
		return err
	})
	return attempt, err
}

func (app *Application) StartExamAttempt(ctx context.Context, examID uuid.UUID, studentID uuid.UUID) (domain.ExamAttempt, error) {
	var attempt domain.ExamAttempt
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		attempt, err = app.examRepo.StartExamAttempt(txCtx, examID, studentID)
		return err
	})
	return attempt, err
}

func (app *Application) SubmitExamAttempt(ctx context.Context, attemptID uuid.UUID, score float64) error {
	return app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return app.examRepo.SubmitExamAttempt(txCtx, attemptID, score)
	})
}

func (app *Application) SaveAnswer(ctx context.Context, attemptID uuid.UUID, questionID uuid.UUID, choiceID uuid.UUID, isCorrect bool) (domain.StudentAnswer, error) {
	var answer domain.StudentAnswer
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		answer, err = app.examRepo.SaveAnswer(txCtx, attemptID, questionID, choiceID, isCorrect)
		return err
	})
	return answer, err
}

func (app *Application) ListAnswersByAttempt(ctx context.Context, attemptID uuid.UUID) ([]domain.StudentAnswer, error) {
	var answers []domain.StudentAnswer
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		answers, err = app.examRepo.ListAnswersByAttempt(txCtx, attemptID)
		return err
	})
	return answers, err
}

func (app *Application) CountExamsBySubject(ctx context.Context, subjectID uuid.UUID) (int64, error) {
	var count int64
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		count, err = app.examRepo.CountExamsBySubject(txCtx, subjectID)
		return err
	})
	return count, err
}

func (app *Application) CountPublishedExamsBySubject(ctx context.Context, subjectID uuid.UUID) (int64, error) {
	var count int64
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		count, err = app.examRepo.CountPublishedBySubject(txCtx, subjectID)
		return err
	})
	return count, err
}

func (app *Application) CountSubmittedAttemptsBySubjectForStudent(ctx context.Context, subjectID uuid.UUID, studentID uuid.UUID) (int64, error) {
	var count int64
	err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		count, err = app.examRepo.CountSubmittedAttemptsBySubjectForStudent(txCtx, subjectID, studentID)
		return err
	})
	return count, err
}

func (app *Application) PublishDraftExamsBySubject(ctx context.Context, subjectID uuid.UUID) error {
	return app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return app.examRepo.PublishDraftExamsBySubject(txCtx, subjectID)
	})
}

func (app *Application) ClosePublishedExamsBySubject(ctx context.Context, subjectID uuid.UUID) error {
	return app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return app.examRepo.ClosePublishedExamsBySubject(txCtx, subjectID)
	})
}

func (app *Application) ListInProgressAttemptsByExam(ctx context.Context, examID uuid.UUID) ([]domain.ExamAttempt, error) {
	var attempts []domain.ExamAttempt
	var err error
	_ = app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		attempts, err = app.examRepo.ListInProgressAttemptsByExam(txCtx, examID)
		return err
	})
	return attempts, err
}

func (app *Application) ListAttemptsByExam(ctx context.Context, examID uuid.UUID) ([]domain.ExamAttempt, error) {
	var attempts []domain.ExamAttempt
	var err error
	_ = app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		attempts, err = app.examRepo.ListAttemptsByExam(txCtx, examID)
		return err
	})
	return attempts, err
}
