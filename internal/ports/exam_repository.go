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
	CountSearch(ctx context.Context, search string) (int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Exam, error)
	Create(ctx context.Context, arg CreateExamParams) (domain.Exam, error)
	Update(ctx context.Context, arg UpdateExamParams) (domain.Exam, error)
	ListBySubject(ctx context.Context, subjectID uuid.UUID) ([]domain.Exam, error)
	SearchBySubject(ctx context.Context, arg SearchExamsBySubjectParams) ([]domain.Exam, error)
	CountSearchBySubject(ctx context.Context, subjectID uuid.UUID, search string) (int64, error)
	ListExamsForStudent(ctx context.Context, studentID uuid.UUID) ([]domain.Exam, error)
	ListExamsCreatedBy(ctx context.Context, instructorID uuid.UUID) ([]domain.Exam, error)
	ListAttemptsByStudent(ctx context.Context, studentID uuid.UUID) ([]domain.ExamAttempt, error)
	GetAttemptByExamAndStudent(ctx context.Context, examID uuid.UUID, studentID uuid.UUID) (domain.ExamAttempt, error)
	StartExamAttempt(ctx context.Context, examID uuid.UUID, studentID uuid.UUID) (domain.ExamAttempt, error)
	SubmitExamAttempt(ctx context.Context, attemptID uuid.UUID, score float64) error
	SaveAnswer(ctx context.Context, attemptID uuid.UUID, questionID uuid.UUID, choiceID uuid.UUID, isCorrect bool) (domain.StudentAnswer, error)
	ListAnswersByAttempt(ctx context.Context, attemptID uuid.UUID) ([]domain.StudentAnswer, error)
	CountExamsBySubject(ctx context.Context, subjectID uuid.UUID) (int64, error)
	CountPublishedBySubject(ctx context.Context, subjectID uuid.UUID) (int64, error)
	CountSubmittedAttemptsBySubjectForStudent(ctx context.Context, subjectID uuid.UUID, studentID uuid.UUID) (int64, error)
	PublishDraftExamsBySubject(ctx context.Context, subjectID uuid.UUID) error
	ClosePublishedExamsBySubject(ctx context.Context, subjectID uuid.UUID) error
	ListInProgressAttemptsByExam(ctx context.Context, examID uuid.UUID) ([]domain.ExamAttempt, error)
	ListAttemptsByExam(ctx context.Context, examID uuid.UUID) ([]domain.ExamAttempt, error)
	ListAttemptsBySubject(ctx context.Context, subjectID uuid.UUID) ([]domain.ExamAttempt, error)
	ListOverallAttemptsBySubject(ctx context.Context, subjectID uuid.UUID) ([]domain.SubjectAttempt, error)
	GetExamAnalytics(ctx context.Context, examID uuid.UUID) (domain.ExamAnalytics, error)
	GetSubjectAnalytics(ctx context.Context, subjectID uuid.UUID) (domain.SubjectAnalytics, error)
}

type SearchExamsParams struct {
	Search string
	Limit  int32
	Offset int32
}

type SearchExamsBySubjectParams struct {
	SubjectID uuid.UUID
	Search    string
	Limit     int32
	Offset    int32
}

type CreateExamParams struct {
	SubjectID   uuid.UUID
	Title       string
	Description *string
	TotalMarks  float64
	PassScore   float64
	StartTime   *string
	EndTime     *string
	Status      domain.ExamStatus
	CreatedBy   uuid.UUID
}

type UpdateExamParams struct {
	ID          uuid.UUID
	Title       string
	Description *string
	TotalMarks  float64
	PassScore   float64
	Status      domain.ExamStatus
}

type ListAllExamsParams struct {
	Limit  int32
	Offset int32
}

type ListAllInstructorsParams struct {
	Limit  int32
	Offset int32
}
