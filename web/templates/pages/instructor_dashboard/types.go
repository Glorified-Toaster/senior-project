package instructor_dashboard

import (
	"uot-exam/internal/domain"
)

type DashboardPageParam struct {
	Instructor domain.User
	Subjects   []SubjectInfo
}

type SubjectInfo struct {
	Subject    domain.Subject
	TotalExams int64
}

type SubjectViewParam struct {
	Instructor  domain.User
	Subject     domain.Subject
	Exams       []domain.Exam
	Students    []domain.User
	AllStudents []domain.User
}
