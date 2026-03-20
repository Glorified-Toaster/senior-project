package repository

import (
	"context"
	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/database/sqlc"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"

	"github.com/google/uuid"
)

type ExamRepository struct {
	queries *sqlc.Queries
}

func NewExamRepository(queries *sqlc.Queries) *ExamRepository {
	return &ExamRepository{queries: queries}
}

func (r *ExamRepository) ListAll(ctx context.Context, arg ports.ListAllExamsParams) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	claimedExams, err := queries.ListAllExams(ctx, sqlc.ListAllExamsParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range claimedExams {
		exams = append(exams, mapSqlcExamToDomain(exam))
	}

	return exams, nil
}

func (r *ExamRepository) Search(ctx context.Context, arg ports.SearchExamsParams) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	claimedExams, err := queries.SearchExams(ctx, sqlc.SearchExamsParams{
		Column1: arg.Search,
		Limit:   arg.Limit,
		Offset:  arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range claimedExams {
		exams = append(exams, mapSqlcExamToDomain(exam))
	}

	return exams, nil
}

func (r *ExamRepository) Count(ctx context.Context) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountExams(ctx)
}

func (r *ExamRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.SoftDeleteExam(ctx, id)
}

func mapSqlcExamToDomain(exam sqlc.Exam) domain.Exam {
	return domain.Exam{
		ID:              exam.ID,
		SubjectID:       exam.SubjectID.UUID,
		Title:           exam.Title,
		Description:     exam.Description,
		DurationMinutes: exam.DurationMinutes,
		TotalMarks:      exam.TotalMarks,
		StartTime:       exam.StartTime.Time,
		EndTime:         exam.EndTime.Time,
		Status:          domain.ExamStatus(exam.Status),
		CreatedBy:       exam.CreatedBy.UUID,
		CreatedAt:       exam.CreatedAt.Time,
		UpdatedAt:       exam.UpdatedAt.Time,
		DeletedAt:       toTimePtr(exam.DeletedAt),
	}
}

func (r *ExamRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	exam, err := queries.GetExamByID(ctx, id)
	if err != nil {
		return domain.Exam{}, err
	}

	return mapSqlcExamToDomain(exam), nil
}
