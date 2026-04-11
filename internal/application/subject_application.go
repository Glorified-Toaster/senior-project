package application

import (
	"context"
	"uot-exam/internal/domain"

	"github.com/google/uuid"
)

func (a *Application) ListAllSubjects(ctx context.Context, limit int32, offset int32) ([]domain.Subject, error) {
	var subjects []domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		subjects, err = a.subjectRepo.ListAllSubjects(txCtx, limit, offset)
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

func (a *Application) SearchSubjects(ctx context.Context, title string, limit int32, offset int32) ([]domain.Subject, error) {
	var subjects []domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		subjects, err = a.subjectRepo.SearchSubjects(txCtx, title, limit, offset)
		return err
	})
	return subjects, err
}

func (a *Application) CountSearchSubjects(ctx context.Context, title string) (int64, error) {
	var count int64
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		count, err = a.subjectRepo.CountSearchSubjects(txCtx, title)
		return err
	})
	return count, err
}

func (a *Application) UpdateSubject(ctx context.Context, subject domain.Subject) (domain.Subject, error) {
	var updatedSubject domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		updatedSubject, err = a.subjectRepo.UpdateSubject(txCtx, subject)
		return err
	})
	return updatedSubject, err
}

func (a *Application) DeleteSubject(ctx context.Context, id uuid.UUID) (domain.Subject, error) {
	var deletedSubject domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		deletedSubject, err = a.subjectRepo.DeleteSubject(txCtx, id)
		return err
	})
	return deletedSubject, err
}

func (a *Application) RestoreSubject(ctx context.Context, id uuid.UUID) (domain.Subject, error) {
	var restoredSubject domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		restoredSubject, err = a.subjectRepo.RestoreSubject(txCtx, id)
		return err
	})
	return restoredSubject, err
}

func (a *Application) CountSubjects(ctx context.Context) (int64, error) {
	var count int64
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		count, err = a.subjectRepo.CountSubjects(txCtx)
		return err
	})
	return count, err
}

func (a *Application) CountDeletedSubjects(ctx context.Context) (int64, error) {
	var count int64
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		count, err = a.subjectRepo.CountDeletedSubjects(txCtx)
		return err
	})
	return count, err
}

func (a *Application) ListDeletedSubjects(ctx context.Context, limit int32, offset int32) ([]domain.Subject, error) {
	var subjects []domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		subjects, err = a.subjectRepo.ListDeletedSubjects(txCtx, limit, offset)
		return err
	})
	return subjects, err
}

func (a *Application) ListInstructorsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]domain.User, error) {
	var instructors []domain.User
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		instructors, err = a.subjectRepo.ListInstructorsBySubjectID(txCtx, subjectID)
		return err
	})
	return instructors, err
}

func (a *Application) CreateSubject(ctx context.Context, subject domain.Subject) (domain.Subject, error) {
	var createdSubject domain.Subject
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		createdSubject, err = a.subjectRepo.CreateSubject(txCtx, subject)
		return err
	})
	return createdSubject, err
}

func (a *Application) DeleteSubjectAndEnrolledInstructors(ctx context.Context, subjectID uuid.UUID) error {
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		err = a.subjectRepo.DeleteSubjectAndEnrolledInstructors(txCtx, subjectID)
		return err
	})
	return err
}

func (a *Application) AssignInstructorToSubject(ctx context.Context, subjectID uuid.UUID, instructorID uuid.UUID) error {
	err := a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		err = a.subjectRepo.AssignInstructorToSubject(txCtx, subjectID, instructorID)
		return err
	})
	return err
}

// AssignInstructorsToSubject assigns each instructor in one transaction; all succeed or none are persisted.
func (a *Application) AssignInstructorsToSubject(ctx context.Context, subjectID uuid.UUID, instructorIDs []uuid.UUID) error {
	if len(instructorIDs) == 0 {
		return nil
	}
	return a.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		for _, id := range instructorIDs {
			if err := a.subjectRepo.AssignInstructorToSubject(txCtx, subjectID, id); err != nil {
				return err
			}
		}
		return nil
	})
}
