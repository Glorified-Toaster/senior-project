package domain

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID           uuid.UUID  `json:"id"`
	ExamID       uuid.UUID  `json:"exam_id"`
	QuestionText string     `json:"question_text"`
	Marks        int        `json:"marks"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
