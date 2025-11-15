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
	PasswordHash   string               `bson:"password_hash" json:"-"`
	IsActive       bool                 `bson:"is_active" json:"is_active"`
	LastLogin      *time.Time           `bson:"last_login,omitempty" json:"last_login,omitempty"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`
	RequiredExams  []primitive.ObjectID `bson:"required_exams,omitempty" json:"required_exams,omitempty"`
	CompletedExams []CompletedExam      `bson:"completed_exams,omitempty" json:"completed_exams,omitempty"`
}

type CompletedExam struct {
	ExamID      primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	Score       float64            `bson:"score" json:"score"`
	TotalMarks  float64            `bson:"total_marks" json:"total_marks"`
	Passed      bool               `bson:"passed" json:"passed"`
	CompletedAt time.Time          `bson:"completed_at" json:"completed_at"`
}
