package response

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentResponse struct {
	ID         primitive.ObjectID `json:"id"`
	FirstName  string             `json:"first_name"`
	LastName   string             `json:"last_name"`
	Role       string             `json:"role"`
	Department string             `json:"department,omitempty"`
	StudentID  string             `json:"student_id"`
	Email      string             `json:"email"`
	IsActive   bool               `json:"is_active"`
	LastLogin  *time.Time         `json:"last_login,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

type PasswordResetResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
