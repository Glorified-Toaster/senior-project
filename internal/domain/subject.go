package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subject struct {
	ID          uuid.UUID
	Title       string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
