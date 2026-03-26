package domain

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID            uuid.UUID  `json:"id"`
	ExamID        uuid.UUID  `json:"exam_id"`
	QuestionTitle string     `json:"question_title"`
	QuestionText  string     `json:"question_text"`
	QuestionType  string     `json:"question_type"`
	QuestionImage string     `json:"question_image"`
	Marks         int        `json:"marks"`
	Choices       []Choice   `json:"choices"`
	Checksum      string     `json:"checksum"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type QuestionType string

const (
	QuestionTypeText  QuestionType = "TEXT"
	QuestionTypeImage QuestionType = "IMAGE"
	QuestionTypeCode  QuestionType = "CODE"
)
