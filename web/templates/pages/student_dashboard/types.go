package student_dashboard

import "time"
import "uot-exam/internal/domain"
import "github.com/google/uuid"

// SubjectInfo holds a subject with its exam count info
type SubjectInfo struct {
	Subject        domain.Subject
	TotalExams     int64
	SubmittedExams int64
}

// ExamInfo holds an exam with its attempt status
type ExamInfo struct {
	Exam    domain.Exam
	Attempt *domain.ExamAttempt
}

// DashboardPageParam holds data for the student dashboard page
type DashboardPageParam struct {
	Student  domain.User
	Subjects []SubjectInfo
}

// SubjectViewParam holds data for the subject detail page
type SubjectViewParam struct {
	Student      domain.User
	Subject      domain.Subject
	Exams        []ExamInfo
	AllSubmitted bool
	ExamEndTime  *time.Time
}

// ExamTakeParam holds data for the exam-taking page
type ExamTakeParam struct {
	Student         domain.User
	Exam            domain.Exam
	Subject         domain.Subject
	Attempt         domain.ExamAttempt
	Questions       []domain.Question
	AnsweredChoices map[uuid.UUID]uuid.UUID
}

// ExamResult holds results for a single exam for the result page
type ExamResult struct {
	Title      string
	Score      int32
	TotalMarks int32
}

// ResultParam holds data for the subject result page
type ResultParam struct {
	Student           domain.User
	Subject           domain.Subject
	ExamResults       []ExamResult
	TotalSubjectScore int32
	MaxPossibleScore  int32
}
