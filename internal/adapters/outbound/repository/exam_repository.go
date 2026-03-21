package repository

import (
	"context"
	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/database/sqlc"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *ExamRepository) Create(ctx context.Context, arg ports.CreateExamParams) (domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	startTime := pgtype.Timestamptz{Valid: false}
	if arg.StartTime != nil && *arg.StartTime != "" {
		t, err := time.Parse("2006-01-02T15:04", *arg.StartTime)
		if err == nil {
			startTime = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}

	endTime := pgtype.Timestamptz{Valid: false}
	if arg.EndTime != nil && *arg.EndTime != "" {
		t, err := time.Parse("2006-01-02T15:04", *arg.EndTime)
		if err == nil {
			endTime = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}

	exam, err := queries.CreateExam(ctx, sqlc.CreateExamParams{
		SubjectID:       uuid.NullUUID{UUID: arg.SubjectID, Valid: true},
		Title:           arg.Title,
		Description:     arg.Description,
		DurationMinutes: arg.DurationMinutes,
		TotalMarks:      arg.TotalMarks,
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          sqlc.ExamStatusType(arg.Status),
		CreatedBy:       uuid.NullUUID{UUID: arg.CreatedBy, Valid: true},
	})
	if err != nil {
		return domain.Exam{}, err
	}

	return mapSqlcExamToDomain(exam), nil
}

func (r *ExamRepository) ListBySubject(ctx context.Context, subjectID uuid.UUID) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcExams, err := queries.ListExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range sqlcExams {
		exams = append(exams, mapSqlcExamToDomain(exam))
	}

	return exams, nil
}
