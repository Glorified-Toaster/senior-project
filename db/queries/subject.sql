-- name: CreateSubject :one
INSERT INTO subjects (title, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetSubjectByID :one
SELECT * FROM subjects WHERE id = $1;

-- name: ListAllSubjects :many
SELECT * FROM subjects WHERE deleted_at IS NULL;

-- name: SearchSubjects :many
SELECT * FROM subjects 
WHERE title LIKE sqlc.arg(title)
AND deleted_at IS NULL 
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: UpdateSubject :one
UPDATE subjects SET title = $2, description = $3, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: DeleteSubject :one
UPDATE subjects SET deleted_at = NOW() WHERE id = $1 RETURNING *;

-- name: RestoreSubject :one
UPDATE subjects SET deleted_at = NULL WHERE id = $1 RETURNING *;

-- name: CountSubjects :one
SELECT COUNT(*) FROM subjects WHERE deleted_at IS NULL;

-- name: CountDeletedSubjects :one
SELECT COUNT(*) FROM subjects WHERE deleted_at IS NOT NULL;

-- name: ListDeletedSubjects :many
SELECT * FROM subjects WHERE deleted_at IS NOT NULL LIMIT $1 OFFSET $2;

-- name: ListInstructorsBySubjectID :many
SELECT 
    u.id,
    u.username,
    u.full_name,
    u.role,
    u.is_active,
    u.last_login,
    u.created_at,
    u.updated_at,
    si.assigned_at,
    u.deleted_at
FROM subject_instructors si
INNER JOIN users u ON si.instructor_id = u.id
WHERE si.subject_id = $1
  AND si.deleted_at IS NULL
  AND u.deleted_at IS NULL
  AND (u.role = 'INSTRUCTOR' OR u.role = 'ADMIN')
ORDER BY u.full_name ASC;

-- name: AssignInstructorToSubject :one
INSERT INTO subject_instructors (subject_id, instructor_id)
VALUES ($1, $2)
RETURNING *;

-- name: UnassignInstructorFromSubject :one
UPDATE subject_instructors SET deleted_at = NOW() WHERE subject_id = $1 AND instructor_id = $2 RETURNING *;