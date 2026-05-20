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
	DeleteChoicesByQuestion(ctx context.Context, arg uuid.UUID) error
	Update(ctx context.Context, arg UpdateQuestionParams) (domain.Question, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Question, error)
	ListAllChoicesByQuestion(ctx context.Context, arg uuid.UUID) ([]domain.Choice, error)
	UpdateChoice(ctx context.Context, arg UpdateChoiceParams) (domain.Choice, error)
	SoftDeleteChoiceByID(ctx context.Context, id uuid.UUID) error
	CountQuestions(ctx context.Context) (int64, error)
}

type UpdateQuestionParams struct {
	ID            uuid.UUID
	QuestionTitle string
	QuestionText  string
	QuestionType  domain.QuestionType
	Marks         float64
	ImageURL      string
	Checksum      string
}

type CreateQuestionParams struct {
	ExamID        uuid.UUID
	QuestionTitle string
	QuestionText  string
	QuestionType  domain.QuestionType
	Marks         float64
	ImageURL      string
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

type UpdateChoiceParams struct {
	ID         uuid.UUID
	ChoiceText string
	IsCorrect  bool
}
