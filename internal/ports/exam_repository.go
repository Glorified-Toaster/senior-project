package ports

import (
	"context"
	"uot-exam/internal/domain"

	"github.com/google/uuid"
)

type ExamRepository interface {
	ListAll(ctx context.Context, arg ListAllExamsParams) ([]domain.Exam, error)
	Search(ctx context.Context, arg SearchExamsParams) ([]domain.Exam, error)
	Count(ctx context.Context) (int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Exam, error)
	Create(ctx context.Context, arg CreateExamParams) (domain.Exam, error)
	ListBySubject(ctx context.Context, subjectID uuid.UUID) ([]domain.Exam, error)
}

type SearchExamsParams struct {
	Search string
	Limit  int32
	Offset int32
}

type CreateExamParams struct {
	SubjectID       uuid.UUID
	Title           string
	Description     *string
	DurationMinutes int32
	TotalMarks      int32
	StartTime       *string
	EndTime         *string
	Status          domain.ExamStatus
	CreatedBy       uuid.UUID
}

type ListAllExamsParams struct {
	Limit  int32
	Offset int32
}

type ListAllInstructorsParams struct {
	Limit  int32
	Offset int32
}
