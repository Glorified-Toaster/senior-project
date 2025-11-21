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

type StudentPasswordHashResponse struct {
	PasswordHash string `bson:"password_hash" json:"-"`
}

type StudentAuthResponse struct {
	StudentResponse *StudentResponse `json:"student"`
	AccessToken     string           `json:"access_token"`
	TokenType       string           `json:"token_type"`
	ExpiresIn       int64            `json:"expires_in"`
}
