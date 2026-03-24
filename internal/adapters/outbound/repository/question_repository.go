package repository

import (
	"context"
	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/database/sqlc"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"

	"github.com/google/uuid"
)

type QuestionRepository struct {
	queries *sqlc.Queries
}

func NewQuestionRepository(queries *sqlc.Queries) *QuestionRepository {
	return &QuestionRepository{queries: queries}
}

func (r *QuestionRepository) Create(ctx context.Context, arg ports.CreateQuestionParams) (domain.Question, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	question, err := queries.CreateQuestion(ctx, sqlc.CreateQuestionParams{
		ExamID:        uuid.NullUUID{UUID: arg.ExamID, Valid: true},
		QuestionTitle: arg.QuestionTitle,
		QuestionText:  arg.QuestionText,
		QuestionType:  sqlc.QuestionTypeType(arg.QuestionType),
		Marks:         int32(arg.Marks),
	})
	if err != nil {
		return domain.Question{}, err
	}

	return mapSqlcQuestionToDomain(question), nil
}

func (r *QuestionRepository) CreateChoice(ctx context.Context, arg ports.CreateChoiceParams) (domain.Choice, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	choice, err := queries.CreateChoice(ctx, sqlc.CreateChoiceParams{
		QuestionID: uuid.NullUUID{UUID: arg.QuestionID, Valid: true},
		ChoiceText: arg.ChoiceText,
		IsCorrect:  arg.IsCorrect,
	})
	if err != nil {
		return domain.Choice{}, err
	}

	return mapSqlcChoiceToDomain(choice), nil
}

func (r *QuestionRepository) ListChoicesByQuestion(ctx context.Context, arg uuid.UUID) ([]domain.Choice, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	choices, err := queries.ListChoicesByQuestion(ctx, uuid.NullUUID{UUID: arg, Valid: true})
	if err != nil {
		return nil, err
	}

	return mapSlice(choices, mapSqlcChoiceToDomain), nil
}

func (r *QuestionRepository) ListQuestionsByExam(ctx context.Context, arg uuid.UUID) ([]domain.Question, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	questions, err := queries.ListQuestionsByExam(ctx, uuid.NullUUID{UUID: arg, Valid: true})
	if err != nil {
		return nil, err
	}

	return mapSlice(questions, mapSqlcQuestionToDomain), nil
}

func mapSqlcQuestionToDomain(question sqlc.Question) domain.Question {
	return domain.Question{
		ID:            question.ID,
		ExamID:        question.ExamID.UUID,
		QuestionTitle: question.QuestionTitle,
		QuestionText:  question.QuestionText,
		QuestionType:  string(question.QuestionType),
		Marks:         int(question.Marks),
		CreatedAt:     question.CreatedAt.Time,
		UpdatedAt:     question.UpdatedAt.Time,
		DeletedAt:     &question.DeletedAt.Time,
	}
}

func mapSlice[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = mapper(v)
	}
	return result
}

func mapSqlcChoiceToDomain(choice sqlc.Choice) domain.Choice {
	return domain.Choice{
		ID:         choice.ID,
		QuestionID: choice.QuestionID.UUID,
		ChoiceText: choice.ChoiceText,
		IsCorrect:  choice.IsCorrect,
		CreatedAt:  choice.CreatedAt.Time,
		UpdatedAt:  choice.UpdatedAt.Time,
		DeletedAt:  &choice.DeletedAt.Time,
	}
}
