-- name: CreateSubject :one
INSERT INTO subjects (title, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetSubjectByID :one
SELECT * FROM subjects WHERE id = $1;

-- name: ListAllSubjects :many
SELECT * FROM subjects WHERE deleted_at IS NULL;

-- name: SearchSubjects :many
SELECT * FROM subjects WHERE title ILIKE '%' || $1::text || '%' AND deleted_at IS NULL;