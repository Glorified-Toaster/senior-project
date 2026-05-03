package domain

import (
	"time"

	"github.com/google/uuid"
)

type ExamStatus string

const (
	ExamStatusDraft     ExamStatus = "DRAFT"
	ExamStatusPublished ExamStatus = "PUBLISHED"
	ExamStatusClosed    ExamStatus = "CLOSED"
)

type Exam struct {
	ID          uuid.UUID
	SubjectID   uuid.UUID
	Title       string
	ExamID      string
	Description *string
	TotalMarks  float64
	PassScore   float64
	StartTime   time.Time
	EndTime     time.Time
	Status      ExamStatus
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time

	TotalQuestions int64
}

type AttemptStatus string

const (
	AttemptStatusInProgress AttemptStatus = "IN_PROGRESS"
	AttemptStatusSubmitted  AttemptStatus = "SUBMITTED"
	AttemptStatusGraded     AttemptStatus = "GRADED"
	AttemptStatusCancelled  AttemptStatus = "CANCELLED"
)

type ExamAttempt struct {
	ID           uuid.UUID
	ExamID       uuid.UUID
	StudentID    uuid.UUID
	StartedAt    time.Time
	SubmittedAt  *time.Time
	Score        *float64
	Status       AttemptStatus
	CreatedAt    time.Time
	ExamTitle       string
	SubjectTitle    string
	StudentName     string
	StudentUsername string
}

type StudentAnswer struct {
	ID               uuid.UUID
	AttemptID        uuid.UUID
	QuestionID       uuid.UUID
	SelectedChoiceID uuid.UUID
	IsCorrect        *bool
	AnsweredAt       time.Time
}

type QuestionAnalytics struct {
	QuestionID     uuid.UUID
	QuestionTitle  string
	QuestionType   string
	MaxMarks       float64
	TotalAnswers   int64
	CorrectAnswers int64
	Accuracy       float64
}

type ExamAnalytics struct {
	ExamID           uuid.UUID
	TotalAttempts    int64
	AverageScore     float64
	MaxScore         float64
	MinScore         float64
	PassCount        int64
	FailCount        int64
	PassRate         float64
	QuestionStats    []QuestionAnalytics
}

type ExamSummaryAnalytics struct {
	ExamID        uuid.UUID
	Title         string
	PassRate      float64
	TotalAttempts int64
}

type SubjectAnalytics struct {
	SubjectID      uuid.UUID
	TotalExams     int64
	TotalAttempts  int64
	AverageScore   float64
	MaxScore       float64
	MinScore       float64
	PassCount      int64
	FailCount      int64
	PassRate       float64
	ExamStats      []ExamSummaryAnalytics
}
