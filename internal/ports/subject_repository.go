package ports

import (
	"context"
	"uot-exam/internal/domain"
)

type SubjectRepository interface {
	ListAllSubjects(ctx context.Context) ([]domain.Subject, error)
	GetSubjectByID(ctx context.Context, id string) (domain.Subject, error)
	SearchSubjects(ctx context.Context, title string) ([]domain.Subject, error)
}

type SearchSubjectsParams struct {
	Title  string
	Limit  int
	Offset int
}
