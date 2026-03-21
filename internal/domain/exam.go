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
	ID              uuid.UUID
	SubjectID       uuid.UUID
	Title           string
	ExamID          string
	Description     *string
	DurationMinutes int32
	TotalMarks      int32
	PassScore       int32
	StartTime       time.Time
	EndTime         time.Time
	Status          ExamStatus
	CreatedBy       uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
