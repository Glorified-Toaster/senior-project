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
	StartTime   time.Time
	EndTime     time.Time
	Status      ExamStatus
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
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
	ExamTitle    string
	SubjectTitle string
}

type StudentAnswer struct {
	ID               uuid.UUID
	AttemptID        uuid.UUID
	QuestionID       uuid.UUID
	SelectedChoiceID uuid.UUID
	IsCorrect        *bool
	AnsweredAt       time.Time
}
