package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubjectStatus string

const (
	SubjectStatusActive    SubjectStatus = "ACTIVE"
	SubjectStatusInactive  SubjectStatus = "INACTIVE"
	SubjectStatusPublished SubjectStatus = "PUBLISHED"
)

type Subject struct {
	ID              uuid.UUID
	Title           string
	Description     *string
	DurationMinutes int32
	TotalMarks      float64
	PassScore       float64

	Status          SubjectStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
