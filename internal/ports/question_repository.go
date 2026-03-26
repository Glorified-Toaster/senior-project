package ports

import (
	"context"
	"uot-exam/internal/domain"

	"github.com/google/uuid"
)

type QuestionRepository interface {
	Create(ctx context.Context, arg CreateQuestionParams) (domain.Question, error)
	CreateChoice(ctx context.Context, arg CreateChoiceParams) (domain.Choice, error)
	ListQuestionsByExam(ctx context.Context, arg uuid.UUID) ([]domain.Question, error)
	ListChoicesByQuestion(ctx context.Context, arg uuid.UUID) ([]domain.Choice, error)
	GetQuestionByChecksum(ctx context.Context, arg string) (bool, error)
	DeleteQuestionAndChoices(ctx context.Context, arg uuid.UUID) error
}

type CreateQuestionParams struct {
	ExamID        uuid.UUID
	QuestionTitle string
	QuestionText  string
	QuestionType  domain.QuestionType
	Marks         int
	Checksum      string
}

type CreateChoiceParams struct {
	QuestionID uuid.UUID
	ChoiceText string
	IsCorrect  bool
}

type ListQuestionsByExamParams struct {
	ExamID uuid.UUID
	Limit  int32
	Offset int32
}

type ListChoicesByQuestionParams struct {
	QuestionID uuid.UUID
	Limit      int32
	Offset     int32
}
