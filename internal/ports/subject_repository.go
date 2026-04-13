package ports

import (
	"context"
	"uot-exam/internal/domain"

	"github.com/google/uuid"
)

type SubjectRepository interface {
	ListAllSubjects(ctx context.Context, limit int32, offset int32) ([]domain.Subject, error)
	GetSubjectByID(ctx context.Context, id string) (domain.Subject, error)
	SearchSubjects(ctx context.Context, title string, limit int32, offset int32) ([]domain.Subject, error)
	CountSearchSubjects(ctx context.Context, title string) (int64, error)
	UpdateSubject(ctx context.Context, subject domain.Subject) (domain.Subject, error)
	DeleteSubject(ctx context.Context, id uuid.UUID) (domain.Subject, error)
	RestoreSubject(ctx context.Context, id uuid.UUID) (domain.Subject, error)
	CountSubjects(ctx context.Context) (int64, error)
	CountDeletedSubjects(ctx context.Context) (int64, error)
	ListDeletedSubjects(ctx context.Context, limit int32, offset int32) ([]domain.Subject, error)
	ListInstructorsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]domain.User, error)
	CreateSubject(ctx context.Context, subject domain.Subject) (domain.Subject, error)
	DeleteSubjectAndEnrolledInstructors(ctx context.Context, subjectID uuid.UUID) error
	AssignInstructorToSubject(ctx context.Context, subjectID uuid.UUID, instructorID uuid.UUID) error
	UnassignInstructorFromSubject(ctx context.Context, subjectID uuid.UUID, instructorID uuid.UUID) error
	AssignStudentsToSubject(ctx context.Context, subjectID uuid.UUID, studentIDs []uuid.UUID) error
	ListStudentsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]domain.User, error)
	ListStudentsBySubjectIDPaginated(ctx context.Context, subjectID uuid.UUID, limit int32, offset int32) ([]domain.User, error)
	CountStudentsBySubjectID(ctx context.Context, subjectID uuid.UUID) (int64, error)
	SearchStudentsBySubjectID(ctx context.Context, subjectID uuid.UUID, search string, limit int32, offset int32) ([]domain.User, error)
	CountSearchStudentsBySubjectID(ctx context.Context, subjectID uuid.UUID, search string) (int64, error)
	UnassignStudentFromSubject(ctx context.Context, subjectID uuid.UUID, studentID uuid.UUID) error
	ListSubjectsForStudent(ctx context.Context, studentID uuid.UUID) ([]domain.Subject, error)
}
