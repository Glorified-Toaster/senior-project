package repository

import (
	"context"
	"uot-exam/internal/adapters/outbound/database/sqlc"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
)

type ExamRepository struct {
	queries *sqlc.Queries
}

func NewExamRepository(queries *sqlc.Queries) *ExamRepository {
	return &ExamRepository{queries: queries}
}

func (r *ExamRepository) ListAll(ctx context.Context, arg ports.ListAllExamsParams) ([]domain.Exam, error) {
	claimedExams, err := r.queries.ListAllExams(ctx, sqlc.ListAllExamsParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range claimedExams {
		exams = append(exams, domain.Exam{
			ID:              exam.ID,
			SubjectID:       exam.SubjectID.UUID,
			Title:           exam.Title,
			Description:     exam.Description,
			DurationMinutes: exam.DurationMinutes,
			TotalMarks:      exam.TotalMarks,
			CreatedBy:       exam.CreatedBy.UUID,
			CreatedAt:       exam.CreatedAt.Time,
			UpdatedAt:       exam.UpdatedAt.Time,
		})
	}

	return exams, nil
}
