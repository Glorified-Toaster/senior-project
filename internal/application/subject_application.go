package application

import (
	"context"
	"uot-exam/internal/domain"

	"github.com/google/uuid"
)

func (a *Application) ListAllSubjects(ctx context.Context) ([]domain.Subject, error) {
	var subjects []domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		subjects, err = a.subjectRepo.ListAllSubjects(txCtx)
		return err
	})
	return subjects, err
}

func (a *Application) GetSubjectByID(ctx context.Context, id uuid.UUID) (domain.Subject, error) {
	var subject domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		subject, err = a.subjectRepo.GetSubjectByID(txCtx, id.String())
		return err
	})
	return subject, err
}

func (a *Application) SearchSubjects(ctx context.Context, title string) ([]domain.Subject, error) {
	var subjects []domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		subjects, err = a.subjectRepo.SearchSubjects(txCtx, title)
		return err
	})
	return subjects, err
}
