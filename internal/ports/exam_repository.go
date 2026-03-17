package ports

import (
	"context"
	"uot-exam/internal/domain"

	"github.com/google/uuid"
)

type ExamRepository interface {
	ListAll(ctx context.Context, arg ListAllExamsParams) ([]domain.Exam, error)
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
