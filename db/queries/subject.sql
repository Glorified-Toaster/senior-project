-- name: CreateSubject :one
INSERT INTO subjects (title, description, duration_minutes, pass_score, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSubjectByID :one
SELECT * FROM subjects WHERE id = $1;

-- name: ListAllSubjects :many
SELECT * FROM subjects 
WHERE deleted_at IS NULL
ORDER BY updated_at DESC, id ASC
LIMIT $1 OFFSET $2;

-- name: SearchSubjects :many
SELECT * FROM subjects 
WHERE title ILIKE '%' || sqlc.arg(title) || '%'
AND deleted_at IS NULL
ORDER BY updated_at DESC, id ASC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountSearchSubjects :one
SELECT COUNT(*)
FROM subjects
WHERE title ILIKE '%' || sqlc.arg(title) || '%'
AND deleted_at IS NULL;

-- name: UpdateSubject :one
UPDATE subjects 
SET title = $2, 
    description = $3, 
    duration_minutes = $4, 
    pass_score = $5, 
    status = $6, 
    updated_at = NOW() 
WHERE id = $1 AND deleted_at IS NULL 
RETURNING *;

-- name: DeleteSubject :one
UPDATE subjects SET deleted_at = NOW() WHERE id = $1 RETURNING *;

-- name: RestoreSubject :one
UPDATE subjects SET deleted_at = NULL WHERE id = $1 RETURNING *;

-- name: CountSubjects :one
SELECT COUNT(*) FROM subjects WHERE deleted_at IS NULL;

-- name: CountDeletedSubjects :one
SELECT COUNT(*) FROM subjects WHERE deleted_at IS NOT NULL;

-- name: ListDeletedSubjects :many
SELECT * FROM subjects 
WHERE deleted_at IS NOT NULL 
ORDER BY deleted_at DESC, id ASC
LIMIT $1 OFFSET $2;

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
  AND u.role = 'INSTRUCTOR'
ORDER BY u.full_name ASC;

-- name: AssignInstructorToSubject :one
INSERT INTO subject_instructors (subject_id, instructor_id)
VALUES ($1, $2)
RETURNING *;

-- name: UnassignInstructorFromSubject :one
UPDATE subject_instructors SET deleted_at = NOW() WHERE subject_id = $1 AND instructor_id = $2 RETURNING *;

-- name: SoftDeleteInstructorsBySubject :exec
UPDATE subject_instructors SET deleted_at = NOW() WHERE subject_id = $1;

-- name: SoftDeleteSubject :exec
UPDATE subjects SET deleted_at = NOW() WHERE id = $1;

-- name: AssignStudentToSubject :one
INSERT INTO subject_students (subject_id, student_id)
VALUES ($1, $2)
RETURNING *;

-- name: UnassignStudentFromSubject :one
UPDATE subject_students SET deleted_at = NOW() WHERE subject_id = $1 AND student_id = $2 RETURNING *;

-- name: ListStudentsBySubjectID :many
SELECT
    u.id,
    u.username,
    u.full_name,
    u.role,
    u.is_active,
    u.last_login,
    u.created_at,
    u.updated_at,
    ss.assigned_at,
    u.deleted_at
FROM subject_students ss
INNER JOIN users u ON ss.student_id = u.id
WHERE ss.subject_id = $1
  AND ss.deleted_at IS NULL
  AND u.deleted_at IS NULL
  AND u.role = 'STUDENT'
ORDER BY u.full_name ASC;

-- name: ListStudentsBySubjectIDPaginated :many
SELECT
    u.id,
    u.username,
    u.full_name,
    u.role,
    u.is_active,
    u.last_login,
    u.created_at,
    u.updated_at,
    ss.assigned_at,
    u.deleted_at
FROM subject_students ss
INNER JOIN users u ON ss.student_id = u.id
WHERE ss.subject_id = $1
  AND ss.deleted_at IS NULL
  AND u.deleted_at IS NULL
  AND u.role = 'STUDENT'
ORDER BY u.full_name ASC
LIMIT $2 OFFSET $3;

-- name: CountStudentsBySubjectID :one
SELECT COUNT(*)
FROM subject_students ss
INNER JOIN users u ON ss.student_id = u.id
WHERE ss.subject_id = $1
  AND ss.deleted_at IS NULL
  AND u.deleted_at IS NULL
  AND u.role = 'STUDENT';

-- name: SearchStudentsBySubjectID :many
SELECT
    u.id,
    u.username,
    u.full_name,
    u.role,
    u.is_active,
    u.last_login,
    u.created_at,
    u.updated_at,
    ss.assigned_at,
    u.deleted_at
FROM subject_students ss
INNER JOIN users u ON ss.student_id = u.id
WHERE ss.subject_id = $1
  AND ss.deleted_at IS NULL
  AND u.deleted_at IS NULL
  AND u.role = 'STUDENT'
  AND (u.full_name ILIKE '%' || $2 || '%' OR u.username ILIKE '%' || $2 || '%')
ORDER BY u.full_name ASC
LIMIT $3 OFFSET $4;

-- name: CountSearchStudentsBySubjectID :one
SELECT COUNT(*)
FROM subject_students ss
INNER JOIN users u ON ss.student_id = u.id
WHERE ss.subject_id = $1
  AND ss.deleted_at IS NULL
  AND u.deleted_at IS NULL
  AND u.role = 'STUDENT'
  AND (u.full_name ILIKE '%' || $2 || '%' OR u.username ILIKE '%' || $2 || '%');

-- name: ListSubjectsForStudent :many
SELECT s.*
FROM subjects s
JOIN subject_students ss ON s.id = ss.subject_id
WHERE ss.student_id = $1
  AND ss.deleted_at IS NULL
  AND s.deleted_at IS NULL
  AND s.status = 'ACTIVE'
ORDER BY s.title ASC;