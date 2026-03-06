package ports

import (
	"context"
	"uot-exam/internal/domain"
)

type ExamRepository interface {
	Create(ctx context.Context, arg CreateExamParams) (domain.Exam, error)
}

type CreateExamParams struct{}
