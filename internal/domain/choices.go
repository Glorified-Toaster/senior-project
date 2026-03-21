package domain

import (
	"time"

	"github.com/google/uuid"
)

type Choice struct {
	ID         uuid.UUID  `json:"id"`
	QuestionID uuid.UUID  `json:"question_id"`
	ChoiceText string     `json:"choice_text"`
	IsCorrect  bool       `json:"is_correct"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
