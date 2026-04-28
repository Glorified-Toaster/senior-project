package repository

import (
	"context"
	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/database/sqlc"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"

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

func (r *ExamRepository) CountSearch(ctx context.Context, search string) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountSearchExams(ctx, search)
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
		ID:          exam.ID,
		SubjectID:   exam.SubjectID.UUID,
		Title:       exam.Title,
		Description: exam.Description,
		TotalMarks:  exam.TotalMarks,
		Status:      domain.ExamStatus(exam.Status),
		CreatedBy:   exam.CreatedBy.UUID,
		CreatedAt:   exam.CreatedAt.Time,
		UpdatedAt:   exam.UpdatedAt.Time,
		DeletedAt:   toTimePtr(exam.DeletedAt),
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

	exam, err := queries.CreateExam(ctx, sqlc.CreateExamParams{
		SubjectID:   uuid.NullUUID{UUID: arg.SubjectID, Valid: true},
		Title:       arg.Title,
		Description: arg.Description,
		TotalMarks:  arg.TotalMarks,
		Status:      sqlc.ExamStatusType(arg.Status),
		CreatedBy:   uuid.NullUUID{UUID: arg.CreatedBy, Valid: true},
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

func (r *ExamRepository) Update(ctx context.Context, arg ports.UpdateExamParams) (domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	exam, err := queries.UpdateExam(ctx, sqlc.UpdateExamParams{
		ID:          arg.ID,
		Title:       arg.Title,
		Description: arg.Description,
		TotalMarks:  arg.TotalMarks,
		Status:      sqlc.ExamStatusType(arg.Status),
	})
	if err != nil {
		return domain.Exam{}, err
	}

	return mapSqlcExamToDomain(exam), nil
}

func (r *ExamRepository) SearchBySubject(ctx context.Context, arg ports.SearchExamsBySubjectParams) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcExams, err := queries.SearchExamsBySubject(ctx, sqlc.SearchExamsBySubjectParams{
		SubjectID: uuid.NullUUID{UUID: arg.SubjectID, Valid: true},
		Column2:   arg.Search,
		Limit:     arg.Limit,
		Offset:    arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range sqlcExams {
		exams = append(exams, mapSqlcExamToDomain(exam))
	}

	return exams, nil
}

func (r *ExamRepository) CountSearchBySubject(ctx context.Context, subjectID uuid.UUID, search string) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountSearchExamsBySubject(ctx, sqlc.CountSearchExamsBySubjectParams{
		SubjectID: uuid.NullUUID{UUID: subjectID, Valid: true},
		Column2:   search,
	})
}

func (r *ExamRepository) ListExamsForStudent(ctx context.Context, studentID uuid.UUID) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcExams, err := queries.ListExamsForStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range sqlcExams {
		exams = append(exams, mapSqlcExamToDomain(exam))
	}

	return exams, nil
}

func (r *ExamRepository) ListExamsCreatedBy(ctx context.Context, instructorID uuid.UUID) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcExams, err := queries.ListExamsCreatedBy(ctx, uuid.NullUUID{UUID: instructorID, Valid: true})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range sqlcExams {
		exams = append(exams, mapSqlcExamToDomain(exam))
	}

	return exams, nil
}

func toInt32Ptr(val pgtype.Int4) *int32 {
	if !val.Valid {
		return nil
	}
	return &val.Int32
}

func mapSqlcAttemptToDomain(attempt sqlc.ExamAttempt) domain.ExamAttempt {
	return domain.ExamAttempt{
		ID:          attempt.ID,
		ExamID:      attempt.ExamID.UUID,
		StudentID:   attempt.StudentID.UUID,
		StartedAt:   attempt.StartedAt.Time,
		SubmittedAt: toTimePtr(attempt.SubmittedAt),
		Score:       toInt32Ptr(attempt.Score),
		Status:      domain.AttemptStatus(attempt.Status),
		CreatedAt:   attempt.CreatedAt.Time,
	}
}

func (r *ExamRepository) ListAttemptsByStudent(ctx context.Context, studentID uuid.UUID) ([]domain.ExamAttempt, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcAttempts, err := queries.ListAttemptsByStudent(ctx, uuid.NullUUID{UUID: studentID, Valid: true})
	if err != nil {
		return nil, err
	}

	var attempts []domain.ExamAttempt
	for _, attempt := range sqlcAttempts {
		attempts = append(attempts, mapSqlcAttemptToDomain(attempt))
	}

	return attempts, nil
}

func (r *ExamRepository) GetAttemptByExamAndStudent(ctx context.Context, examID uuid.UUID, studentID uuid.UUID) (domain.ExamAttempt, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	attempt, err := queries.GetAttemptByExamAndStudent(ctx, sqlc.GetAttemptByExamAndStudentParams{
		ExamID:    uuid.NullUUID{UUID: examID, Valid: true},
		StudentID: uuid.NullUUID{UUID: studentID, Valid: true},
	})
	if err != nil {
		return domain.ExamAttempt{}, err
	}

	return mapSqlcAttemptToDomain(attempt), nil
}

func (r *ExamRepository) StartExamAttempt(ctx context.Context, examID uuid.UUID, studentID uuid.UUID) (domain.ExamAttempt, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	attempt, err := queries.StartExamAttempt(ctx, sqlc.StartExamAttemptParams{
		ExamID:    uuid.NullUUID{UUID: examID, Valid: true},
		StudentID: uuid.NullUUID{UUID: studentID, Valid: true},
	})
	if err != nil {
		return domain.ExamAttempt{}, err
	}

	return mapSqlcAttemptToDomain(attempt), nil
}

func (r *ExamRepository) SubmitExamAttempt(ctx context.Context, attemptID uuid.UUID, score int32) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.SubmitExamAttempt(ctx, sqlc.SubmitExamAttemptParams{
		ID:    attemptID,
		Score: pgtype.Int4{Int32: score, Valid: true},
	})
}

func mapSqlcAnswerToDomain(answer sqlc.StudentAnswer) domain.StudentAnswer {
	var isCorrect *bool
	if answer.IsCorrect.Valid {
		isCorrect = &answer.IsCorrect.Bool
	}
	return domain.StudentAnswer{
		ID:               answer.ID,
		AttemptID:        answer.AttemptID.UUID,
		QuestionID:       answer.QuestionID.UUID,
		SelectedChoiceID: answer.SelectedChoiceID.UUID,
		IsCorrect:        isCorrect,
		AnsweredAt:       answer.AnsweredAt.Time,
	}
}

func (r *ExamRepository) SaveAnswer(ctx context.Context, attemptID uuid.UUID, questionID uuid.UUID, choiceID uuid.UUID, isCorrect bool) (domain.StudentAnswer, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	answer, err := queries.SaveAnswer(ctx, sqlc.SaveAnswerParams{
		AttemptID:        uuid.NullUUID{UUID: attemptID, Valid: true},
		QuestionID:       uuid.NullUUID{UUID: questionID, Valid: true},
		SelectedChoiceID: uuid.NullUUID{UUID: choiceID, Valid: true},
		IsCorrect:        pgtype.Bool{Bool: isCorrect, Valid: true},
	})
	if err != nil {
		return domain.StudentAnswer{}, err
	}

	return mapSqlcAnswerToDomain(answer), nil
}

func (r *ExamRepository) ListAnswersByAttempt(ctx context.Context, attemptID uuid.UUID) ([]domain.StudentAnswer, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcAnswers, err := queries.ListAnswersByAttempt(ctx, uuid.NullUUID{UUID: attemptID, Valid: true})
	if err != nil {
		return nil, err
	}

	var answers []domain.StudentAnswer
	for _, answer := range sqlcAnswers {
		answers = append(answers, mapSqlcAnswerToDomain(answer))
	}

	return answers, nil
}

func (r *ExamRepository) CountExamsBySubject(ctx context.Context, subjectID uuid.UUID) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
}

func (r *ExamRepository) CountPublishedBySubject(ctx context.Context, subjectID uuid.UUID) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountPublishedExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
}

func (r *ExamRepository) CountSubmittedAttemptsBySubjectForStudent(ctx context.Context, subjectID uuid.UUID, studentID uuid.UUID) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountSubmittedAttemptsBySubjectForStudent(ctx, sqlc.CountSubmittedAttemptsBySubjectForStudentParams{
		SubjectID: uuid.NullUUID{UUID: subjectID, Valid: true},
		StudentID: uuid.NullUUID{UUID: studentID, Valid: true},
	})
}

func (r *ExamRepository) PublishDraftExamsBySubject(ctx context.Context, subjectID uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.PublishDraftExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
}
