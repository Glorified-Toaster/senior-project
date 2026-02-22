-- name: CreateSubject :one
INSERT INTO subjects (title, description, instructor_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSubjectByID :one
SELECT * FROM subjects WHERE id = $1;

-- name: ListSubjectsByInstructor :many
SELECT * FROM subjects WHERE instructor_id = $1 ORDER BY created_at DESC;


-- name: EnrollStudent :one
INSERT INTO enrollments (student_id, subject_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetEnrollmentsByStudent :many
SELECT * FROM enrollments WHERE student_id = $1;

-- name: GetStudentsInSubject :many
SELECT student_id FROM enrollments WHERE subject_id = $1;
