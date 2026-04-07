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

func (r *SubjectRepository) ListAllSubjects(ctx context.Context, limit int32, offset int32) ([]domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	subjects, err := queries.ListAllSubjects(ctx, sqlc.ListAllSubjectsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var domainSubjects []domain.Subject
	for _, subject := range subjects {
		domainSubjects = append(domainSubjects, domain.Subject{
			ID:              subject.ID,
			Title:           subject.Title,
			Description:     subject.Description,
			DurationMinutes: subject.DurationMinutes,
			TotalMarks:      subject.TotalMarks,
			PassScore:       subject.PassScore,
			Status:          domain.SubjectStatus(subject.Status),
			CreatedAt:       subject.CreatedAt.Time,
			UpdatedAt:       subject.UpdatedAt.Time,
			DeletedAt:       toTimePtr(subject.DeletedAt),
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
		ID:              subject.ID,
		Title:           subject.Title,
		Description:     subject.Description,
		DurationMinutes: subject.DurationMinutes,
		TotalMarks:      subject.TotalMarks,
		PassScore:       subject.PassScore,
		Status:          domain.SubjectStatus(subject.Status),
		CreatedAt:       subject.CreatedAt.Time,
		UpdatedAt:       subject.UpdatedAt.Time,
		DeletedAt:       toTimePtr(subject.DeletedAt),
	}, nil
}

func (r *SubjectRepository) SearchSubjects(ctx context.Context, title string, limit int32, offset int32) ([]domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	subjects, err := queries.SearchSubjects(ctx, sqlc.SearchSubjectsParams{
		Title:      &title,
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, err
	}

	var domainSubjects []domain.Subject
	for _, subject := range subjects {
		domainSubjects = append(domainSubjects, domain.Subject{
			ID:              subject.ID,
			Title:           subject.Title,
			Description:     subject.Description,
			DurationMinutes: subject.DurationMinutes,
			TotalMarks:      subject.TotalMarks,
			PassScore:       subject.PassScore,
			Status:          domain.SubjectStatus(subject.Status),
			CreatedAt:       subject.CreatedAt.Time,
			UpdatedAt:       subject.UpdatedAt.Time,
			DeletedAt:       toTimePtr(subject.DeletedAt),
		})
	}

	return domainSubjects, nil
}

func (r *SubjectRepository) CountSearchSubjects(ctx context.Context, title string) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	// Wildcard matching is handled in the SQL query itself.
	count, err := queries.CountSearchSubjects(ctx, &title)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *SubjectRepository) UpdateSubject(ctx context.Context, subject domain.Subject) (domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	updatedSubject, err := queries.UpdateSubject(ctx, sqlc.UpdateSubjectParams{
		ID:              subject.ID,
		Title:           subject.Title,
		Description:     subject.Description,
		DurationMinutes: subject.DurationMinutes,
		TotalMarks:      subject.TotalMarks,
		PassScore:       subject.PassScore,
		Status:          sqlc.SubjectStatusType(subject.Status),
	})
	if err != nil {
		return domain.Subject{}, err
	}

	return domain.Subject{
		ID:              updatedSubject.ID,
		Title:           updatedSubject.Title,
		Description:     updatedSubject.Description,
		DurationMinutes: updatedSubject.DurationMinutes,
		TotalMarks:      updatedSubject.TotalMarks,
		PassScore:       updatedSubject.PassScore,
		Status:          domain.SubjectStatus(updatedSubject.Status),
		CreatedAt:       updatedSubject.CreatedAt.Time,
		UpdatedAt:       updatedSubject.UpdatedAt.Time,
		DeletedAt:       toTimePtr(updatedSubject.DeletedAt),
	}, nil
}

func (r *SubjectRepository) DeleteSubject(ctx context.Context, id uuid.UUID) (domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	deletedSubject, err := queries.DeleteSubject(ctx, id)
	if err != nil {
		return domain.Subject{}, err
	}

	return domain.Subject{
		ID:              deletedSubject.ID,
		Title:           deletedSubject.Title,
		Description:     deletedSubject.Description,
		DurationMinutes: deletedSubject.DurationMinutes,
		TotalMarks:      deletedSubject.TotalMarks,
		PassScore:       deletedSubject.PassScore,
		Status:          domain.SubjectStatus(deletedSubject.Status),
		CreatedAt:       deletedSubject.CreatedAt.Time,
		UpdatedAt:       deletedSubject.UpdatedAt.Time,
		DeletedAt:       toTimePtr(deletedSubject.DeletedAt),
	}, nil
}

func (r *SubjectRepository) RestoreSubject(ctx context.Context, id uuid.UUID) (domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	restoredSubject, err := queries.RestoreSubject(ctx, id)
	if err != nil {
		return domain.Subject{}, err
	}

	return domain.Subject{
		ID:              restoredSubject.ID,
		Title:           restoredSubject.Title,
		Description:     restoredSubject.Description,
		DurationMinutes: restoredSubject.DurationMinutes,
		TotalMarks:      restoredSubject.TotalMarks,
		PassScore:       restoredSubject.PassScore,
		Status:          domain.SubjectStatus(restoredSubject.Status),
		CreatedAt:       restoredSubject.CreatedAt.Time,
		UpdatedAt:       restoredSubject.UpdatedAt.Time,
		DeletedAt:       toTimePtr(restoredSubject.DeletedAt),
	}, nil
}

func (r *SubjectRepository) CountSubjects(ctx context.Context) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	count, err := queries.CountSubjects(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SubjectRepository) CountDeletedSubjects(ctx context.Context) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	count, err := queries.CountDeletedSubjects(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SubjectRepository) ListDeletedSubjects(ctx context.Context, limit int32, offset int32) ([]domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	subjects, err := queries.ListDeletedSubjects(ctx, sqlc.ListDeletedSubjectsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var domainSubjects []domain.Subject
	for _, subject := range subjects {
		domainSubjects = append(domainSubjects, domain.Subject{
			ID:              subject.ID,
			Title:           subject.Title,
			Description:     subject.Description,
			DurationMinutes: subject.DurationMinutes,
			TotalMarks:      subject.TotalMarks,
			PassScore:       subject.PassScore,
			Status:          domain.SubjectStatus(subject.Status),
			CreatedAt:       subject.CreatedAt.Time,
			UpdatedAt:       subject.UpdatedAt.Time,
			DeletedAt:       toTimePtr(subject.DeletedAt),
		})
	}

	return domainSubjects, nil
}

func (r *SubjectRepository) ListInstructorsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]domain.User, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	instructors, err := queries.ListInstructorsBySubjectID(ctx, subjectID)
	if err != nil {
		return nil, err
	}

	var domainInstructors []domain.User
	for _, instructor := range instructors {
		domainInstructors = append(domainInstructors, domain.User{
			ID:        instructor.ID,
			Username:  instructor.Username,
			FullName:  instructor.FullName,
			Role:      domain.UserRole(instructor.Role),
			IsActive:  instructor.IsActive,
			LastLogin: toTimePtr(instructor.LastLogin),
			CreatedAt: instructor.CreatedAt.Time,
			UpdatedAt: instructor.UpdatedAt.Time,
			DeletedAt: toTimePtr(instructor.DeletedAt),
		})
	}

	return domainInstructors, nil
}

func (r *SubjectRepository) CreateSubject(ctx context.Context, subject domain.Subject) (domain.Subject, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	createdSubject, err := queries.CreateSubject(ctx, sqlc.CreateSubjectParams{
		Title:           subject.Title,
		Description:     subject.Description,
		DurationMinutes: subject.DurationMinutes,
		TotalMarks:      subject.TotalMarks,
		PassScore:       subject.PassScore,
		Status:          sqlc.SubjectStatusType(subject.Status),
	})
	if err != nil {
		return domain.Subject{}, err
	}

	return domain.Subject{
		ID:          createdSubject.ID,
		Title:       createdSubject.Title,
		Description: createdSubject.Description,
		CreatedAt:   createdSubject.CreatedAt.Time,
		UpdatedAt:   createdSubject.UpdatedAt.Time,
		DeletedAt:   toTimePtr(createdSubject.DeletedAt),
	}, nil
}

func (r *SubjectRepository) DeleteSubjectAndEnrolledInstructors(ctx context.Context, subjectID uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	err := queries.SoftDeleteInstructorsBySubject(ctx, subjectID)
	if err != nil {
		return err
	}

	err = queries.SoftDeleteSubject(ctx, subjectID)
	if err != nil {
		return err
	}

	return nil
}
