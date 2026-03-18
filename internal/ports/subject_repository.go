package ports

import (
	"context"
	"uot-exam/internal/domain"
)

type SubjectRepository interface {
	ListAllSubjects(ctx context.Context) ([]domain.Subject, error)
	GetSubjectByID(ctx context.Context, id string) (domain.Subject, error)
	SearchSubjects(ctx context.Context, title string) ([]domain.Subject, error)
	UpdateSubject(ctx context.Context, subject domain.Subject) (domain.Subject, error)
	DeleteSubject(ctx context.Context, id string) (domain.Subject, error)
	RestoreSubject(ctx context.Context, id string) (domain.Subject, error)
	CountSubjects(ctx context.Context) (int64, error)
	CountDeletedSubjects(ctx context.Context) (int64, error)
	ListDeletedSubjects(ctx context.Context, limit int32, offset int32) ([]domain.Subject, error)
	ListInstructorsBySubjectID(ctx context.Context, subjectID string) ([]domain.User, error)
}

type SearchSubjectsParams struct {
	Title  string
	Limit  int
	Offset int
}
