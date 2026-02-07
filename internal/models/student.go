package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	FirstName      string               `bson:"first_name" json:"first_name" validate:"required,min=2,max=32"`
	LastName       string               `bson:"last_name" json:"last_name" validate:"required,min=2,max=32"`
	Role           string               `bson:"role" json:"role"`
	Department     string               `bson:"department,omitempty" json:"department,omitempty"`
	StudentID      string               `bson:"student_id" json:"student_id"`
	Email          string               `bson:"email" json:"email" validate:"email,required"`
	PasswordHash   string               `bson:"password_hash" json:"password_hash,omitempty"`
	IsActive       bool                 `bson:"is_active" json:"is_active"`
	LastLogin      *time.Time           `bson:"last_login,omitempty" json:"last_login,omitempty"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`
	RequiredExams  []primitive.ObjectID `bson:"required_exams,omitempty" json:"required_exams,omitempty"`
	CompletedExams []CompletedExam      `bson:"completed_exams,omitempty" json:"completed_exams,omitempty"`
}

func (s *Student) GetID() string       { return s.StudentID }
func (s *Student) GetRole() string     { return "student" }
func (s *Student) GetEmail() string    { return s.Email }
func (s *Student) GetPassword() string { return s.PasswordHash }
func (s *Student) IsActiveUser() bool  { return s.IsActive }
