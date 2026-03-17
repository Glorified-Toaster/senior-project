package repository

import (
	"context"

	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/database/sqlc"
	"uot-exam/internal/domain"

	"github.com/google/uuid"
)

type SubjectRepository struct {
	queries *sqlc.Queries
}

func NewSubjectRepository(queries *sqlc.Queries) *SubjectRepository {
	return &SubjectRepository{queries: queries}
}

func (r *SubjectRepository) ListAllSubjects(ctx context.Context) ([]domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	subjects, err := queries.ListAllSubjects(ctx)
	if err != nil {
		return nil, err
	}

	var domainSubjects []domain.Subject
	for _, subject := range subjects {
		domainSubjects = append(domainSubjects, domain.Subject{
			ID:          subject.ID,
			Title:       subject.Title,
			Description: subject.Description,
			CreatedAt:   subject.CreatedAt.Time,
			UpdatedAt:   subject.UpdatedAt.Time,
			DeletedAt:   toTimePtr(subject.DeletedAt),
		})
	}

	return domainSubjects, nil
}

func (r *SubjectRepository) GetSubjectByID(ctx context.Context, id string) (domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	subject, err := queries.GetSubjectByID(ctx, uuid.MustParse(id))
	if err != nil {
		return domain.Subject{}, err
	}

	return domain.Subject{
		ID:          subject.ID,
		Title:       subject.Title,
		Description: subject.Description,
		CreatedAt:   subject.CreatedAt.Time,
		UpdatedAt:   subject.UpdatedAt.Time,
		DeletedAt:   toTimePtr(subject.DeletedAt),
	}, nil
}

func (r *SubjectRepository) SearchSubjects(ctx context.Context, title string) ([]domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	subjects, err := queries.SearchSubjects(ctx, title)
	if err != nil {
		return nil, err
	}

	var domainSubjects []domain.Subject
	for _, subject := range subjects {
		domainSubjects = append(domainSubjects, domain.Subject{
			ID:          subject.ID,
			Title:       subject.Title,
			Description: subject.Description,
			CreatedAt:   subject.CreatedAt.Time,
			UpdatedAt:   subject.UpdatedAt.Time,
			DeletedAt:   toTimePtr(subject.DeletedAt),
		})
	}

	return domainSubjects, nil
}
