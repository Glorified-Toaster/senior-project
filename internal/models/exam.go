package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Exam Struct
type CompletedExam struct {
	ExamID      primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	Score       float64            `bson:"score" json:"score"`
	TotalMarks  float64            `bson:"total_marks" json:"total_marks"`
	Passed      bool               `bson:"passed" json:"passed"`
	CompletedAt time.Time          `bson:"completed_at" json:"completed_at"`
}
